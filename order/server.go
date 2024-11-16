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
	"google.golang.org/grpc/reflection"
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

// PostOrder processes a new order request, validating the account and products before creation
func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	// Verify account exists by calling account service
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("Error getting account:", err)
		return nil, errors.New("account not found")
	}

	// Extract product IDs from the request for catalog lookup
	var ids []string
	for _, p := range r.Products {
		ids = append(ids, p.ProductId)
	}

	// Fetch full product details from catalog service
	products, err := s.catalogClient.GetProducts(ctx, 0, 0, ids, "")
	if err != nil {
		log.Println("Error getting products:", err)
		return nil, errors.New("products not found")
	}

	// Prepare ordered products by combining catalog details with requested quantities
	var orderedProducts []OrderedProduct
	for _, p := range products {
		// Initialize product with catalog details
		product := OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0,
		}

		// Find matching product from request to get quantity
		for _, pro := range r.Products {
			if pro.ProductId == p.ID {
				product.Quantity = pro.Quantity
				break
			}
		}

		// Only include products that were actually ordered (quantity > 0)
		if product.Quantity != 0 {
			orderedProducts = append(orderedProducts, product)
		}
	}

	// Create the order in the order service
	order, err := s.service.PostOrder(ctx, r.AccountId, orderedProducts)
	if err != nil {
		log.Println("Error posting order: ", err)
		return nil, errors.New("could not post order")
	}

	// Convert domain order to protobuf format
	orderProto := &pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
	}

	// Convert timestamp to binary format for protobuf
	orderProto.CreatedAt, err = order.CreatedAt.MarshalBinary()
	if err != nil {
		log.Println("Error conveting CreatedAt time to byte: ", err)
		return nil, errors.New("couldn't convert time to byte CreateAt")
	}

	// Add ordered products to protobuf response
	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}

	// Return the created order
	return &pb.PostOrderResponse{Order: orderProto}, nil
}

// GetOrdersForAccount retrieves orders for a specific account and enriches them with product details from the catalog service
func (s *grpcServer) GetOrdersForAccount(ctx context.Context, r *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	// Get all orders for the account from the order service
	accOrders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("Error getting account orders: ", err)
		return nil, err
	}

	// Create a map to deduplicate product IDs across all orders
	productIDMap := map[string]bool{}
	for _, o := range accOrders {
		for _, p := range o.Products {
			productIDMap[p.ID] = true
		}
	}

	// Convert the map keys (product IDs) to a slice for the catalog service call
	productIDs := []string{}
	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}

	// Fetch product details from the catalog service
	products, err := s.catalogClient.GetProducts(ctx, 1, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting products with ids from catalog")
		return nil, err
	}

	// Convert domain orders to protobuf orders and enrich with product details
	orders := []*pb.Order{}
	for _, o := range accOrders {
		// Create new protobuf order
		op := &pb.Order{
			Id:         o.ID,
			AccountId:  o.AccountID,
			TotalPrice: o.TotalPrice,
		}

		// Convert time.Time to binary for protobuf
		op.CreatedAt, err = o.CreatedAt.MarshalBinary()
		if err != nil {
			log.Println("Error conveting time to bytes: ", err)
			return nil, err
		}

		// Enrich each ordered product with details from catalog
		for _, product := range o.Products {
			// Find matching product from catalog and update details
			for _, p := range products {
				if p.ID == product.ID {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price
					break
				}
			}

			// Add enriched product to protobuf order
			op.Products = append(op.Products, &pb.Order_OrderProduct{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    product.Quantity,
			})
		}

		// Add the completed order (with all its enriched products) to the final orders slice
		orders = append(orders, op)
	}

	// Return the response with all enriched orders
	return &pb.GetOrdersForAccountResponse{Orders: orders}, nil
}
