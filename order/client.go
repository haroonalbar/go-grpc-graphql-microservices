package order

import (
	"context"
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

func (c *Client) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	protoProducts := []*pb.PostOrderRequest_OrderProduct{}
	for _, p := range products {
		protoProducts = append(protoProducts, &pb.PostOrderRequest_OrderProduct{
			ProductId: p.ID,
			Quantity:  p.Quantity,
		})
	}

	res, err := c.service.PostOrder(ctx, &pb.PostOrderRequest{
		AccountId: accountID,
		Products:  protoProducts,
	})
	if err != nil {
		return nil, err
	}

	newOrder := &Order{
		ID:         res.Order.Id,
		AccountID:  res.Order.AccountId,
		TotalPrice: res.Order.TotalPrice,
		Products:   products,
	}

	err = newOrder.CreatedAt.UnmarshalBinary(res.Order.CreatedAt)
	if err != nil {
		return nil, err
	}

	return newOrder, nil
}

func (c *Client) GetOrderForAccount(ctx context.Context, accountID string) ([]Order, error) {
	res, err := c.service.GetOrdersForAccount(ctx, &pb.GetOrdersForAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	orders := []Order{}

	for _, orderProto := range res.Orders {
		newOrder := Order{
			ID:         orderProto.Id,
			AccountID:  orderProto.AccountId,
			TotalPrice: orderProto.TotalPrice,
		}

		err := newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)
		if err != nil {
			log.Println("Error while converting bytes to time : ", err)
			return nil, err
		}

		products := []OrderedProduct{}
		for _, p := range orderProto.Products {
			products = append(products, OrderedProduct{
				ID:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    p.Quantity,
			})
		}
		newOrder.Products = products

		orders = append(orders, newOrder)
	}

	return orders, nil
}
