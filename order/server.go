package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/haroonalbar/go-grpc-graphql-microservices/account"
	"github.com/haroonalbar/go-grpc-graphql-microservices/catalog"
	"github.com/haroonalbar/go-grpc-graphql-microservices/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client // Dependency: Account microservice client
	catalogClient *catalog.Client // Dependency: Catalog microservice client
}

// ListenGRPC starts the gRPC server, establishing connections to Account and Catalog services.
// It also registers the OrderService server and handles graceful cleanup of resources.
func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {
	// Attempt to connect to the Account service
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return fmt.Errorf("failed to connect to account service: %w", err)
	}
	defer accountClient.Close() // Ensures cleanup if initialization fails

	// Attempt to connect to the Catalog service
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		return fmt.Errorf("failed to connect to catalog service: %w", err)
	}
	defer catalogClient.Close() // Ensures cleanup if initialization fails

	// Start listening on the specified TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to start listener on port %d: %w", port, err)
	}

	// Create a new gRPC server
	serv := grpc.NewServer()

	// Register OrderService with gRPC server
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	})

	// Register reflection service for debugging (consider restricting this in production)
	reflection.Register(serv)

	// Defer server closure for graceful shutdown in production scenarios
	defer func() {
		serv.GracefulStop()
	}()

	// Start serving requests
	return serv.Serve(lis)
}

// PostOrder processes a new order request.
func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	// get account
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("Error getting account:", err)
		return nil, errors.New("account not found")
	}

	// get product ids from request
	var ids []string
	for _, p := range r.Products {
		ids = append(ids, p.ProductId)
	}
	// get products from catalog
	products, err := s.catalogClient.GetProducts(ctx, 0, 0, ids, "")
	if err != nil {
		log.Println("Error getting products:", err)
		return nil, errors.New("products not found")
	}

	// for calling PostOrder
	var orderedProducts []OrderedProduct
	for _, p := range products {
		product := OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0,
		}
		for _, pro := range r.Products {
			if pro.ProductId == p.ID {
				product.Quantity = pro.Quantity
				break
			}
		}
		if product.Quantity != 0 {
			orderedProducts = append(orderedProducts, product)
		}
	}

	// get order
	order, err := s.service.PostOrder(ctx, r.AccountId, orderedProducts)
	if err != nil {
		log.Println("Error posting order: ", err)
		return nil, errors.New("could not post order")
	}

	orderProto := &pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
	}

	// convert CreatedAt to bytes
	orderProto.CreatedAt, err = order.CreatedAt.MarshalBinary()
	if err != nil {
		log.Println("Error conveting CreatedAt time to byte: ", err)
		return nil, errors("couldn't convert time to byte CreateAt")
	}

	// add products
	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}
	return &pb.PostOrderResponse{Order: orderProto}, nil
}

// GetOrdersForAccount retrieves orders for a specific account.
// This is a placeholder implementation and currently returns Unimplemented.
func (s *grpcServer) GetOrdersForAccount(ctx context.Context, req *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	// Return an Unimplemented error to indicate this method is a placeholder.
	return nil, status.Error(codes.Unimplemented, "method GetOrdersForAccount not implemented")
}
