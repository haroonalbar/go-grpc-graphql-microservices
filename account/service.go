package account

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Account struct{
	ID string `json:"id"`
	Name string `json:"name"`
}

// The Service layer acts as an intermediary between the Server (API handlers) and the Repository (data access).
// It encapsulates the business logic of the application, providing several benefits:
//
// 1. Separation of Concerns: It separates the business logic from the data access and presentation layers,
//    making the codebase more modular and easier to maintain.
// 2. Abstraction: It provides a clean API for the Server layer, hiding the complexities of data operations
//    and allowing for easier mocking in tests.
// 3. Business Logic: It's where we implement domain-specific logic, data validation, and orchestration
//    of multiple repository calls if needed.
// 4. Scalability: As the application grows, having a separate Service layer makes it easier to add
//    new features or modify existing ones without affecting other parts of the application.
// 5. Reusability: Business logic in the Service layer can be reused across different parts of the application
//    or even in different applications.
// 6. Transaction Management: If needed, the Service layer can manage database transactions that span
//    multiple repository operations.
//
// By using this layered architecture (Server -> Service -> Repository), we create a more robust,
// maintainable, and scalable application structure.

type Service interface{
	PostAccount (ctx context.Context, name string) (*Account, error)
	GetAccount (ctx context.Context, id string) (*Account, error)
	GetAccounts (ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type accountService struct{
	repository Repository
}


func(s *accountService) PostAccount (ctx context.Context, name string) (*Account, error){
	// creates new id using ksuid
	a := Account{ID: ksuid.New().String(), Name: name}
	
	if err := s.repository.PutAccount(ctx, a);err != nil{
		return nil , err
	}
	return &a, nil
}

func (s *accountService) GetAccount (ctx context.Context, id string) (*Account, error){
	return s.repository.GetAccountByID(ctx, id)
}

func (s *accountService) GetAccounts (ctx context.Context, skip uint64, take uint64) ([]Account, error){
	if take > 100 || (skip == 0 && take == 0){
		take = 100
	}
	return s.repository.ListAccounts(ctx,skip,take)
}

func newService(r Repository) Service{
	return &accountService{r}
}