package order

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrderForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgresRepository and attempts to open a connection to the database using
// the given connection string. If the connection attempt fails, it returns an error.
func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

// PutOrder saves an order to the database.
//
// It does this by using a single database transaction to:
// 1. Insert the main order record into the orders table.
// 2. Bulk insert all the products associated with the order, using the PostgreSQL COPY command.
//
// The function returns an error if any part of this process fails.
// The key benefits of this implementation are:
// Atomicity: Either all operations succeed or none do
// Performance: Uses COPY for efficient bulk insertion of products
// Safety: Proper
// error handling and transaction management
// Context Support: Respects context cancellation
//
// The deferred function should check if err is not nil before deciding to commit or rollback,
// but it should use a named return error to capture the function's scope error correctly
func (r *postgresRepository) PutOrder(ctx context.Context, o Order) (err error) {
	// Starts a new database transaction. All subsequent operations will be part of this atomic transaction.
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	// Deferred Transaction Management
	defer func() {
		// 		This ensures that if any error occurs during the transaction:
		// The transaction is rolled back (all changes are undone)
		// If successful, the transaction is committed
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	// Insert Order
	// Inserts the main order record into the orders table.
	_, err = tx.ExecContext(ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES ($1, $2, $3, $4)",
		o.ID,
		o.CreatedAt,
		o.AccountID,
		o.TotalPrice,
	)
	if err != nil {
		return
	}
	// Bulk Insert Products
	// Uses PostgreSQL's COPY command (through pq.CopyIn) for efficient bulk insertion of order products.
	stmt, err := tx.PrepareContext(ctx,
		pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	if err != nil {
		return
	}
	// Insert Each Product
	// Loops through each product and adds it to the bulk insert operation.
	for _, p := range o.Products {
		_, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity)
		if err != nil {
			return
		}
	}
	// Finalize Bulk Insert:
	// 	Executes the bulk insert of all products.
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return
	}

	stmt.Close()
	return
}

// This function is using a row-by-row processing approach to group products
// with their respective orders, taking advantage of the ORDER BY o.id to ensure
// all products for the same order are processed together.

// GetOrderForAccount retrieves all orders for the given account ID.
//
// It executes a single SQL query to retrieve all orders and their products
// in a single pass. This query uses a JOIN to combine the orders and order_products
// tables. The result is then processed in a single pass, using a row-by-row
// approach to group products with their respective orders. The result is a slice
// of Order objects, each containing a slice of OrderedProduct objects.
func (r *postgresRepository) GetOrderForAccount(ctx context.Context, accountID string) ([]Order, error) {
	// Execute SQL query to get orders and their products
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
		o.id,
		o.created_at,
		o.account_id,
		o.total_price::money::numeric::float8,  // Convert money to float8
		op.product_id,
		op.quantity
		FROM orders o JOIN order_products op ON(o.id = op.order_id)
		WHERE o.account_id = $1
		ORDER BY o.id
		`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize tracking variables
	currentOrder := &Order{}            // Holds the order currently being processed
	lastOrder := &Order{}               // Holds the previous order for comparison
	orders := []Order{}                 // Final slice of all orders
	orderedProduct := &OrderedProduct{} // Temporary holder for product data
	products := []OrderedProduct{}      // Collects products for current order

	// Iterate through result rows
	for rows.Next() {
		// Scan current row into our structs
		if err := rows.Scan(
			&currentOrder.ID,
			&currentOrder.CreatedAt,
			&currentOrder.AccountID,
			&currentOrder.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}

		// If we've moved to a new order (ID changed)
		if lastOrder.ID != "" && lastOrder.ID != currentOrder.ID {
			// Create and append the completed order
			newOrder := Order{
				ID:         lastOrder.ID,
				AccountID:  lastOrder.AccountID,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products:   products, // Assign collected products
			}
			orders = append(orders, newOrder)

			// Reset products slice for new order
			products = []OrderedProduct{}
		}

		// Add current product to products slice
		products = append(products, OrderedProduct{
			ID:       orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})

		// Update lastOrder for next iteration
		*lastOrder = *currentOrder
	}

	// Handle the last order after loop ends
	if lastOrder.ID != "" {
		newOrder := Order{
			ID:         lastOrder.ID,
			AccountID:  lastOrder.AccountID,
			CreatedAt:  lastOrder.CreatedAt,
			TotalPrice: lastOrder.TotalPrice,
			Products:   products, // Assign collected products
		}
		orders = append(orders, newOrder)
	}

	// Check for any errors that occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
