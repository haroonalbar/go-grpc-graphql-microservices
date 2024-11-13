package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/haroonalbar/go-grpc-graphql-microservices/order"
)

var ErrInvalidParameter = errors.New("invalid parameter")

// mutationResolver resolves Mutation. ie, create, update, delete
type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, in AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	acc, err := r.server.accountClient.PostAccount(ctx, in.Name)
	if err != nil {
		log.Println("Error creating account from account client: ", err)
		return nil, err
	}
	return &Account{
		ID:   acc.ID,
		Name: acc.Name,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTiemout(ctx, 3*time.Second)
	defer cancel()

	p, err := r.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.Price)
	if err != nil {
		log.Println("Error creating product on catalog client: ", err)
		return nil, err
	}

	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	products := []order.OrderedProduct
	for _, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		products = append(products, p)
	}

	o, err := r.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil {
		log.Println("Error creating order in order client: ", err)
		return nil, err
	}

	return &Order{
		ID:         o.ID,
		CreatedAt:  o.CreatedAt,
		TotalPrice: o.TotalPrice,
		Products:   o.Products,
	}, nil
}
