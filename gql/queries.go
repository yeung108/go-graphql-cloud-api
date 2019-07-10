package gql

import (
	"context"
	"fmt"

	"go-graphql-cloud-api/postgres"
	"log"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// Root holds a pointer to a graphql object
type Root struct {
	Query   *graphql.Object
	Context *context.Context
}

func GetVendorProductsBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		var result dataloader.Result
		result.Error = err
		results = append(results, &result)
		return results
	}
	var vendorIDs []uuid.UUID
	for _, key := range keys {
		k, err := uuid.FromString(key.String())
		if err != nil {
			fmt.Println("vendorIDs key error: %v ", err)
		}
		vendorIDs = append(vendorIDs, k)
	}
	products, err := keys[0].(*ResolverKey).client().resolver().db.GetVendorProducts(vendorIDs)
	if err != nil {
		return handleError(err)
	}

	var results []*dataloader.Result
	result := dataloader.Result{
		Data:  products,
		Error: nil,
	}
	results = append(results, &result)

	log.Printf("[GetVendorProductsBatchFn] batch size: %d", len(results))
	return results
}

func GetVendorStoresBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		var result dataloader.Result
		result.Error = err
		results = append(results, &result)
		return results
	}
	var vendorIDs []uuid.UUID
	for _, key := range keys {
		k, err := uuid.FromString(key.String())
		if err != nil {
			fmt.Println("vendorIDs key error: %v ", err)
		}
		vendorIDs = append(vendorIDs, k)
	}
	stores, err := keys[0].(*ResolverKey).client().resolver().db.GetVendorStores(vendorIDs)
	if err != nil {
		return handleError(err)
	}

	var results []*dataloader.Result
	result := dataloader.Result{
		Data:  stores,
		Error: nil,
	}
	results = append(results, &result)

	log.Printf("[GetVendorStoresBatchFn] batch size: %d", len(results))
	return results
}

func GetVendorsBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		var result dataloader.Result
		result.Error = err
		results = append(results, &result)
		return results
	}
	var vendorIDs []uuid.UUID
	for _, key := range keys {
		k, err := uuid.FromString(key.String())
		if err != nil {
			fmt.Println("vendorIDs key error: %v ", err)
		}
		vendorIDs = append(vendorIDs, k)
	}
	vendors, err := keys[0].(*ResolverKey).client().resolver().db.GetVendors(vendorIDs)
	if err != nil {
		return handleError(err)
	}

	var results []*dataloader.Result
	result := dataloader.Result{
		Data:  vendors,
		Error: nil,
	}
	results = append(results, &result)
	log.Printf("[GetVendorsBatchFn] batch size: %d", len(results))
	return results
}

// NewRoot returns base query type. This is where we add all the base queries
func NewRoot(db *postgres.Db) *Root {
	// Create a resolver holding our databse. Resolver can be found in resolvers.go
	resolver := Resolver{db: db}
	// Dataloader
	var loaders = make(map[string]*dataloader.Loader, 1)
	var client = Client{Resolver: &resolver}
	loaders["GetVendorProducts"] = dataloader.NewBatchedLoader(GetVendorProductsBatchFn)
	loaders["GetVendorStores"] = dataloader.NewBatchedLoader(GetVendorStoresBatchFn)
	loaders["GetVendors"] = dataloader.NewBatchedLoader(GetVendorsBatchFn)
	ctx := context.WithValue(context.Background(), "loaders", loaders)
	ctx = context.WithValue(ctx, "client", &client)

	// Create a new Root that describes our base query set up. In this
	// example we have a user query that takes one argument called name
	root := Root{
		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "Query",
				Fields: graphql.Fields{
					"vendors": &graphql.Field{
						// Slice of Vendor type which can be found in types.go
						Type: graphql.NewList(Vendor),
						Args: graphql.FieldConfigArgument{
							"id": &graphql.ArgumentConfig{
								Type: graphql.String,
							},
						},
						Resolve: resolver.VendorResolver,
					},
				},
			},
		),
		Context: &ctx,
	}
	return &root
}
