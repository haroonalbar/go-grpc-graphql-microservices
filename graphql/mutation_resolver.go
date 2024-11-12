package main

import "context"

// mutationResolver resolves Mutation. ie, create, update, delete
type mutationResolver struct {
	server *Server
}

// CreateAccount(ctx context.Context, account AccountInput) (*Account, error)
// CreateProduct(ctx context.Context, product ProductInput) (*Product, error)
// CreateOrder(ctx context.Context, order OrderInput) (*Order, error)
func (r *mutationResolver) CreateAccount(ctx context.Context, in AccountInput) (*Account, error) {
	panic("")
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	panic("")
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	panic("")
}
