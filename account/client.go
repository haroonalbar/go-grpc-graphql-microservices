package account

import (
	"context"

	"github.com/haroonalbar/go-grpc-graphql-microservices/account/pb"
	"google.golang.org/grpc"
)

// NOTE:
// This will used from the graphql/Server struct
// As per the flow graphql to client to server to service to repository to db
// So this Client is used by graphql to intact with account Microservice
type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

// returns a Client with grpc connection and account service client
func NewClient(url string) (*Client, error) {
	// making a grpc connection
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	// getting account service client from connection
	c := pb.NewAccountServiceClient(conn)
	return &Client{conn, c}, nil
}

// close connection
func (c *Client) Close() {
	c.conn.Close()
}

// Calling all of it from the pb generated file
func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.PostAccount(
		ctx,
		&pb.PostAccountRequest{
			Name: name,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccount(
		ctx,
		&pb.GetAccountRequest{
			Id: id,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	r, err := c.service.GetAccounts(
		ctx,
		&pb.GetAccountsRequest{
			Skip: skip,
			Take: take,
		},
	)
	if err != nil {
		return nil, err
	}
	var accounts []Account
	for _, acc := range r.Accounts {
		accounts = append(accounts, Account{
			ID:   acc.Id,
			Name: acc.Name,
		})
	}
	return accounts, nil
}
