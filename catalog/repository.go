package catalog

import (
	"context"
	"errors"

	elastic "github.com/olivere/elastic/v7"
)

type Product struct{}

// WARN: using depricated elastic search client
// update the implementation to the official client
// "github.com/elastic/go-elasticsearch/v8"

var ErrNotFound = errors.New("Entity not found")

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

func (r *elasticRepository) Close() {}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	panic("")
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	panic("")
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	panic("")
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	panic("")
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	panic("")
}

type productDocument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	// WHY:is price string?
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		// "sniffing" in the context of Elasticsearch client libraries,
		// it refers to the ability of these clients to dynamically discover and connect to nodes in an Elasticsearch cluster.
		// This feature helps clients maintain connections to the cluster even if individual nodes change or become unavailable.
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}

	return &elasticRepository{
		client: client,
	}, nil
}
