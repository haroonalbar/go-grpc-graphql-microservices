//go:generate protoc --go_out=. --go-grpc_out=. account.proto
package account

import (
	"context"
	"fmt"
	"net"

	"github.com/haroonalbar/go-grpc-graphql-microservices/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is going to call funcs in Service
// Also going to communicate with Client using gRPC protobuff data packets

// This struct embeds a Service interface, which likely defines the actual business logic for account operations.
type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

// TODO:
// fix error on RegisterAccountServiceServer

// This function sets up and starts the gRPC server
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
	pb.RegisterAccountServiceServer(serv, &grpcServer{service: s})
	// Starts serving gRPC requests.
	return serv.Serve(lis)
}

// pb is generated from protobuff file
// here *pb. will be PostAccountRequest and Response
func (s *grpcServer) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	a, err := s.service.PostAccount(ctx, r.Name)
	if err != nil {
		return nil, err
	}
	return &pb.PostAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	a, err := s.service.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	res, err := s.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	accounts := []*pb.Account{}
	for _, a := range res {
		accounts = append(accounts, &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		})
	}
	return &pb.GetAccountsResponse{
		Accounts: accounts,
	}, nil
}
