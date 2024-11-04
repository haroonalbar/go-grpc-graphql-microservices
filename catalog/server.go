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
	p, err := s.service.PostProduct(ctx, r.Name, r.Description, r.Price)
	if err != nil {
		return nil, err
	}

	return &pb.PostProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	ps, err := s.service.GetProducts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	var products []*pb.Product
	for _, p := range ps {
		products = append(products, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return &pb.GetProductsResponse{
		Products: products,
	}, nil
}
