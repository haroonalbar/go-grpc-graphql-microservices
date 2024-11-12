package order

import (
	"log"

	"github.com/haroonalbar/go-grpc-graphql-microservices/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	// NOTE: used NewClient instead of depricated Dial
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := pb.NewOrderServiceClient(conn)
	return &Client{conn: conn, service: c}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Println("Error closing grpc connection : ", err)
	}
}

func (c *Client) PostOrder() {
	panic("")
}

func (c *Client) GetOrderForAccount() {
	panic("")
}
