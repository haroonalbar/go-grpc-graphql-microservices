package order

import (
	"context"
	"fmt"
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
// This is a placeholder implementation and currently returns Unimplemented.
func (s *grpcServer) PostOrder(ctx context.Context, req *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	// Return an Unimplemented error to indicate this method is a placeholder.
	return nil, status.Error(codes.Unimplemented, "method PostOrder not implemented")
}

// GetOrdersForAccount retrieves orders for a specific account.
// This is a placeholder implementation and currently returns Unimplemented.
func (s *grpcServer) GetOrdersForAccount(ctx context.Context, req *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	// Return an Unimplemented error to indicate this method is a placeholder.
	return nil, status.Error(codes.Unimplemented, "method GetOrdersForAccount not implemented")
}
