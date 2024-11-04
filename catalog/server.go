//go:generate protoc --go_out=. --go-grpc_out=. catalog.proto
package catalog

import (
	"context"
	"fmt"
	"net"

	"github.com/haroonalbar/go-grpc-graphql-microservices/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	// It creates a TCP listener on the specified port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	// Initializes a new gRPC server.
	serv := grpc.NewServer()
	// Registers the server for reflection (useful for debugging and service discovery).
	reflection.Register(serv)

	// Register the service
	pb.RegisterCatalogServiceServer(
		serv,
		&grpcServer{
			service:                           s,
			UnimplementedCatalogServiceServer: pb.UnimplementedCatalogServiceServer{},
		},
	)

	// Starts serving gRPC requests
	return serv.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	panic("")
}

func (s *grpcServer) GetProduct(context.Context, *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	panic("")
}

func (s *grpcServer) GetProducts(context.Context, *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	panic("")
}
