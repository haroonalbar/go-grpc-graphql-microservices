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
	conn *grpc.ClientConn
	// from generated pb not from catalog
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	// conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	res, err := c.service.PostProduct(
		ctx,
		&pb.PostProductRequest{
			Name:        name,
			Description: description,
			Price:       price,
		})
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	res, err := c.service.GetProduct(
		ctx,
		&pb.GetProductRequest{
			Id: id,
		})
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, skip, take uint64, ids []string, query string) ([]Product, error) {
	res, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Ids:   ids,
		Skip:  skip,
		Take:  take,
		Query: query,
	})
	if err != nil {
		return nil, err
	}
	var products []Product
	for _, p := range res.Products {
		products = append(products, Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}
