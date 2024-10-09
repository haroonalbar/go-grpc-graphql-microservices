package main

import "github.com/99designs/gqlgen/graphql"

type Server struct {
	// accountClient *account.Client
	// catalogClient *catalog.Client
	// orderClient   *order.Client
}

func NewGraphQLServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {
	// accountClient, err := account.NewClient(accountUrl)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// // catalogClient is dependant on accountClient
	// catalogClient, err := catalog.NewClient(catalogUrl)
	// if err != nil {
	// 	accountClient.Close()
	// 	return nil, err
	// }
	//
	// // orderClient is dependant on both clients above
	// orderClient, err := order.NewClient(orderUrl)
	// if err != nil {
	// 	accountClient.Close()
	// 	catalogClient.Close()
	// 	return nil, err
	// }
	//
	// return &Server{
	// 	accountClient,
	// 	catalogClient,
	// 	orderClient,
	// }, nil
	return nil, nil
}

// func (s *Server) Mutation() MutationResolver {
// 	// mutation_resolver.go
// 	return &mutationResolver{
// 		server: s,
// 	}
// }
//
// func (s *Server) Query() QueryResolver {
// 	// query_resolver.go
// 	return &queryResolver{
// 		server: s,
// 	}
// }
//
// func (s *Server) Account() AccountResolver {
// 	// account_resolver.go
// 	return &accountResolver{
// 		server: s,
// 	}
// }

// // ToExecutableSchema is converting the Server instance into an executable GraphQL schema.
func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	// return NewExecutableSchema(Config{
	// 	// server matches ResolverRoot interface in graphql
	// 	Resolvers: s,
	// })
	panic("")
}
