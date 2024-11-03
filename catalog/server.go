//go:generate protoc --go_out=. --go-grpc_out=. catalog.proto
package catalog

import "github.com/haroonalbar/go-grpc-graphql-microservices/catalog/pb"

type grpcServer struct {
	pb.UnimplementedProductServiceServer
	service Service
}
