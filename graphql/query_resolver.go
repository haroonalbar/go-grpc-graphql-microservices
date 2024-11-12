package main

import "context"

type queryResolver struct {
	server *Server
}

// Orders is in accounts_resolver cause it's depandant on account

// Accounts takes in pagination and id cause that's the signature we defined in graphql.schema
// and generated by gqlgen
func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	panic("")
}

// Products samething as Accounts with extra parameter query that's also defined in schema
func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	panic("")
}
