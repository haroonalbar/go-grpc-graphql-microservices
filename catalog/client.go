package catalog

import (
	"context"

	"github.com/haroonalbar/go-grpc-graphql-microservices/catalog/pb"
	"google.golang.org/grpc"
)

// NOTE: This will used from the graphql/Server struct
// As per the flow graphql to client to server to service to repository to db
// So this Client is used by graphql to intact with catalog Microservice
type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := pb.NewCatalogServiceClient(conn)

	return &Client{
		conn:    conn,
		service: c,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

