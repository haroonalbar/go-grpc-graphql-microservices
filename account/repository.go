package account

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq" // postgres driver
)

// The Repository interface and its implementation (postgresRepository) serve several important purposes:
// The Repository pattern is a design pattern used to abstract the ** Data access layer ** from the business logic layer (service) in an application.
//
// 1. Separation of Concerns: It isolates the data access logic from the rest of the application,
//    making it easier to maintain and modify the data layer without affecting other parts of the system.
// 2. Abstraction: It provides a clean, abstract interface for data operations, hiding the complexities
//    of database interactions from the rest of the application.
// 3. Testability: By abstracting the data access, it becomes easier to mock the repository for unit testing
//    other parts of the application that depend on data access.
// 4. Flexibility: It allows for easy switching between different data sources (e.g., from a relational
//    database to a NoSQL database) without changing the application's business logic.
// 5. Centralized Data Logic: It provides a centralized place to implement data access logic, promoting
//    code reuse and reducing duplication.
// 6. Type Safety: By defining specific methods for data operations, it provides type safety and
//    reduces the likelihood of runtime errors due to incorrect SQL queries.
// 7. Performance Optimization: It allows for the implementation of caching strategies and query optimizations
//    in a centralized location.
// 8. Transactional Control: It provides a natural place to implement transaction management for operations
//    that span multiple database calls.
//
// By implementing these patterns and principles, we create a more modular, maintainable, and scalable
// application architecture that can easily adapt to changing requirements and technologies.

type Repository interface {
	Close()
	PutAccount(ctx context.Context, a Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) PutAccount(ctx context.Context, a Account) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts(id, name) VALUES($1, $2)", a.ID, a.Name)
	return err
}

func (r *postgresRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name FROM accounts WHERE id = $1", id)
	a := &Account{}
	if err := row.Scan(&a.ID, &a.Name); err != nil {
		return nil, err
	}
	return a, nil
}

// ListAccounts returns a slice of Account objects, sorted by ID descending,
// with 'skip' number of elements skipped and 'take' number of elements returned.
// The method returns an error if the query fails.
func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2",
		skip,
		take,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	accounts:= []Account{}
	for rows.Next(){
		a:= &Account{}
		if err = rows.Scan(&a.ID, &a.Name); err == nil {
			accounts = append(accounts, *a)
		}
	}
	if err = rows.Err(); err != nil{
		return nil, err
	}
	return accounts,nil
}

func NewPostgresRepository(url string) (Repository, error) {
	// connect to db
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	// check if the connection is established
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}
