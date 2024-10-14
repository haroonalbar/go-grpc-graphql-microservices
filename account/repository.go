package account

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq" // postgres driver
)

type Account struct {
	ID   string
	Name string
}

// The Repository pattern is a design pattern used to abstract the data access layer from the business logic layer in an application.
// By using this pattern, we create a more robust, maintainable, and scalable application architecture.
//
// The Repository interface and its implementation (postgresRepository) serve several important purposes:

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
