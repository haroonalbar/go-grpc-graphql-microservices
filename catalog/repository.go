package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"gopkg.in/olivere/elastic.v5"
)

type productDocument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	// WARN: why is price string?
}

type Product struct {
	ID          string
	Name        string
	Description string
	Price       string
}

// TODO: using depricated elastic search client
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
	client *elasticsearch.Client
	// // depricated
	// renamed to clientdep for continuation with the vid
	clientdep *elastic.Client
}

func (r *elasticRepository) Close() {}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	jsonBody, err := json.Marshal(productDocument{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	})
	if err != nil {
		return fmt.Errorf("Failed to marshal product: %w", err)
	}

	// Index creates or updates a document in an index.
	res, err := r.client.Index(
		"catalog",
		bytes.NewReader(jsonBody),
		// r.client.Index.WithDocumentType("product"),
		r.client.Index.WithDocumentID(p.ID),
		r.client.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("Failed to execute index request: %w", err)
	}
	defer res.Body.Close()

	// Check if response idicates an err
	if res.IsError() {
		return fmt.Errorf("Error response form Elasticsearch: %w", err)
	}
	return nil

	// // depricated
	// // new index is stored at index catalog of type product . set id to p.id
	// // and set body of document from productDocument as fields and formatted as json
	// // and execute the request through Do() ctx used for controlling the lifetime of the req
	// _, err := r.client.Index().
	// 	Index("catalog").
	// 	Type("product").
	// 	Id(p.ID).
	// 	BodyJson(productDocument{
	// 		Name:        p.Name,
	// 		Description: p.Description,
	// 		Price:       p.Price,
	// 	}).Do(ctx)
	// return err
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	// official
	res, err := r.client.Get(
		"catalog",
		id,
		r.client.Get.WithContext(ctx),
		// r.client.Get.WithDocumentType("porduct"),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to get product: %w", err)
	}
	defer res.Body.Close()

	// Check if document was not found
	if res.StatusCode == 404 {
		return nil, ErrNotFound
	}

	// Check for other errors
	if res.IsError() {
		return nil, fmt.Errorf("error getting document: %s", res.String())
	}

	// Parse the response body
	// struct that matches Elasticserch response structure
	var response struct {
		Source productDocument `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("Error parsing response body: %w", err)
	}

	return &Product{
		ID:          id,
		Name:        response.Source.Name,
		Description: response.Source.Description,
		Price:       response.Source.Price,
	}, nil

	// // Another way
	// // Parse the response body
	// var result map[string]interface{}
	// if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
	// 	return nil, fmt.Errorf("Error parsing response body: %w", err)
	// }
	// // Extract the _source field
	// source, ok := result["_source"].(map[string]interface{})
	// if !ok {
	// 	return nil, fmt.Errorf("_source not found in response")
	// }
	// return &Product{
	// 	ID:          id,
	// 	Name:        source["name"].(string),
	// 	Description: source["description"].(string),
	// 	Price:       source["price"].(string),
	// }, nil

	// // Depricated
	// res, err := r.client.Get().Index("catalog").Type("product").Id(id).Do(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// // if not found
	// if !res.Found {
	// 	return nil, ErrNotFound
	// }
	// p := productDocument{}
	// if err := json.Unmarshal(*res.Source, &p); err != nil {
	// 	return nil, err
	// }
	// return &Product{
	// 	ID:          id,
	// 	Name:        p.Name,
	// 	Description: p.Description,
	// 	Price:       p.Price,
	// }, nil
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {

	// // Create the search query
	// // SearchRequest configures the Search API request.
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"match_all": map[string]interface{}{},
	// 	},
	// 	"from": int(skip),
	// 	"size": int(take),
	// }

	// // Convert query to JSON
	// var buf bytes.Buffer
	// if err := json.NewEncoder(&buf).Encode(query); err != nil {
	// 	return nil, fmt.Errorf("Error encoding query: %w", err)
	// }

	// Perform the search request
	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithSearchType("product"),
		// r.client.Search.WithBody(&buf),
		// WARN: Can i use the below code instead of teh WithBody ??
		r.client.Search.WithFrom(int(skip)),
		r.client.Search.WithSize(int(take)),
		r.client.Search.WithQuery("match_all"),
	)
	if err != nil {
		return nil, fmt.Errorf("Error while performing search: %w", err)
	}
	defer res.Body.Close()

	// Parse the response
	var result struct {
		Hits struct {
			Hits []struct {
				ID     string          `json:"_id"`
				Source productDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Error parsing response: %w", err)
	}

	// Convert hits to products
	products := make([]Product, 0, len(result.Hits.Hits))

	for _, hit := range result.Hits.Hits {
		products = append(products, Product{
			ID:          hit.ID,
			Name:        hit.Source.Name,
			Description: hit.Source.Description,
			Price:       hit.Source.Price,
		})
	}

	return products, nil

	// // Depricated
	// res, err := r.clientdep.Search().
	// 	Index("catalog")
	// 	Type("product").
	// 	Query(elastic.NewMatchAllQuery()).
	// 	From(int(skip)).
	// 	Size(int(take)).Do(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// products := []Product{}
	// for _, hit := range res.Hits.Hits {
	// 	p := productDocument{}
	// 	if err = json.Unmarshal(*hit.Source, &p); err == nil {
	// 		products = append(products, Product{
	// 			ID:          hit.Id,
	// 			Name:        p.Name,
	// 			Description: p.Description,
	// 			Price:       p.Price,
	// 		})
	// 	}
	// }
	//
	// return products, nil
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {

	// Prepare the request body docs
	docs := make([]map[string]string, len(ids))
	for i, id := range ids {
		docs[i] = map[string]string{"_id": id}
	}

	// // convert body to JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(map[string]interface{}{"ids": ids}); err != nil {
		return nil, fmt.Errorf("error encoding body: %w", err)
	}

	// Perform the mget request
	res, err := r.client.Mget(
		&buf,
		r.client.Mget.WithIndex("catalog"),
		r.client.Mget.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("error performing mget: %w", err)
	}
	defer res.Body.Close()

	// Parse the response
	var result struct {
		Docs []struct {
			ID     string          `json:"_id"`
			Source productDocument `json:"_source"`
			Found  bool            `json:"found"`
		} `json:"docs"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	// convert results to  products
	products := make([]Product, 0, len(result.Docs))
	for _, doc := range result.Docs {
		if doc.Found {
			products = append(products, Product{
				ID:          doc.ID,
				Name:        doc.Source.Name,
				Description: doc.Source.Description,
				Price:       doc.Source.Price,
			})
		}
	}

	return products, nil

	// // unofficial depricated way
	// items := []*elastic.MultiGetItem{}
	// for _, id := range ids {
	// 	items = append(
	// 		items,
	// 		elastic.NewMultiGetItem().
	// 			Index("catalog").
	// 			Type("product").
	// 			Id(id),
	// 	)
	// }

	// resp, err := r.clientdep.
	// 	MultiGet().
	// 	Add(items...).
	// 	Do(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error calling muti get: %w", err)
	// }

	// var products []Product

	// for _, doc := range resp.Docs {
	// 	var p productDocument
	// 	if err := json.Unmarshal(*doc.Source, &p); err != nil {
	// 		return nil, fmt.Errorf("Error deconding response: %w", err)
	// 	} else {
	// 		products = append(products, Product{
	// 			ID:          doc.Id,
	// 			Name:        p.Name,
	// 			Description: p.Description,
	// 			Price:       p.Price,
	// 		})
	// 	}
	// }

	// return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	panic("")
}

func NewElasticRepository(url string) (Repository, error) {
	// official
	// by default will use port 9200
	// and [http.DefaultTransport]
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}

	res, err := es.Info()
	if err != nil {
		return nil, err
	}

	fmt.Print(res)

	// return &elasticRepository{
	// 	client: es,

	// }, nil

	// not official
	client, err := elastic.NewClient(
		// SetURL defines the URL endpoints of the Elasticsearch nodes.
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
		client:    es,
		clientdep: client,
	}, nil
}
