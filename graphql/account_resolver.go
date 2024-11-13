package main

import (
	"context"
	"log"
	"time"
)

type accountResolver struct {
	server *Server
}

// Orders are in accountResolver cause it's dependent on account
func (r *accountResolver) Orders(ctx context.Context, acc *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := r.server.orderClient.GetOrderForAccount(ctx, acc.ID)
	if err != nil {
		log.Println("Error getting orders for account from order client: ", err)
		return nil, err
	}
	var orders []*Order
	for _, o := range orderList {
		var products []*OrderedProduct
		for _, op := range o.Products {
			products = append(products, &OrderedProduct{
				ID:          op.ID,
				Name:        op.Name,
				Description: op.Description,
				Price:       op.Price,
				Quantity:    int(op.Quantity),
			})
		}
		orders = append(orders, &Order{
			ID:         o.ID,
			CreatedAt:  o.CreatedAt,
			TotalPrice: o.TotalPrice,
			Products:   products,
		})
	}
	return orders, nil
}
