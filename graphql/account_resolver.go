package main

import "context"

type accountResolver struct {
	server *Server
}

// Orders are in accountResolver cause it's dependent on account
func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	panic("")
}
