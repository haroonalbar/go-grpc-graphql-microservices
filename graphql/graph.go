package main

// type Server struct {
// 	accountClient *account.Client
// 	catalogClient *catalog.Client
// 	orderClient   *order.Client
// }

func NewGraphQLServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}

	// catalogClient is dependant on accountClient
	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	// orderClient is dependant on both clients above
	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, err
	}

  return &Server{
    accountClient
    catalogClient
    orderClient
  }, nil
}
