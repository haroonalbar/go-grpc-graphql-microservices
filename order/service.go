package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type Order struct {
	ID         string           `json:"id"`
	AccountID  string           `json:"account_id"`
	CreatedAt  time.Time        `json:"created_at"`
	TotalPrice float64          `json:"tatal_price"`
	Products   []OrderedProduct `json:"products"`
}

type OrderedProduct struct {
	ID          string  `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    uint32  `json:"quantity"`
}

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type orderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &orderService{r}
}

func (s orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	o := Order{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountID: accountID,
		Products:  products,
	}
	for _, p := range products {
		// by default it would be 0.0
		o.TotalPrice += p.Price * float64(p.Quantity)
	}
	err := s.repository.PutOrder(ctx, o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (s orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return s.repository.GetOrderForAccount(ctx, accountID)
}
