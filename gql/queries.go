package gql

import (
	"context"
	"errors"
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

	var vendorsMap = make(map[uuid.UUID][]postgres.Product, len(vendorIDs))
	for _, product := range products {
		if _, found := vendorsMap[product.VendorID.UUID]; found {
			vendorsMap[product.VendorID.UUID] = append(
				vendorsMap[product.VendorID.UUID], product)
		} else {
			vendorsMap[product.VendorID.UUID] = []postgres.Product{product}
		}
	}

	var results []*dataloader.Result
	for _, vendorID := range vendorIDs {
		products, ok := vendorsMap[vendorID]
		if !ok {
			err := errors.New(fmt.Sprintf("products not found, "+
				"vendor id: %d", vendorID))
			return handleError(err)
		}
		result := dataloader.Result{
			Data:  products,
			Error: nil,
		}
		results = append(results, &result)
	}
	log.Printf("[GetVendorProductsBatchFn] batch size: %d", len(results))
	return results
}

// NewRoot returns base query type. This is where we add all the base queries
func NewRoot(db *postgres.Db) *Root {
	// Create a resolver holding our databse. Resolver can be found in resolvers.go
	resolver := Resolver{db: db}
	newID, err := uuid.FromString("fbc77caf-7dfb-46b4-883f-786b98e012e5")
	if err != nil {
		fmt.Println("Initial Parsing Error when parsing %s", err)
	}
	fmt.Println("newID: %+v", newID)
	currentVendor := postgres.Vendor{
		ID: newID,
	}
	// Dataloader TODO
	var loaders = make(map[string]*dataloader.Loader, 1)
	var client = Client{Resolver: &resolver}
	loaders["GetVendorProducts"] = dataloader.NewBatchedLoader(GetVendorProductsBatchFn)
	ctx := context.WithValue(context.Background(), "currentVendor", &currentVendor)
	ctx = context.WithValue(ctx, "loaders", loaders)
	ctx = context.WithValue(ctx, "client", &client)

	// Create a new Root that describes our base query set up. In this
	// example we have a user query that takes one argument called name
	root := Root{
		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "Query",
				Fields: graphql.Fields{
					"currentVendor": &graphql.Field{
						Type: Vendor,
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							var vendor = p.Context.Value("currentVendor").(*postgres.Vendor)
							return vendor, nil
						},
					},
					"vendors": &graphql.Field{
						// Slice of Setting type which can be found in types.go
						Type: graphql.NewList(Vendor),
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							var (
								v       = p.Context.Value
								c       = v("client").(*Client)
								loaders = v("loaders").(map[string]*dataloader.Loader)
								vendor  = p.Source.(*postgres.Vendor)
								key     = NewResolverKey(fmt.Sprintf("%d", vendor.ID), c)
							)
							thunk := loaders["GetVendorProducts"].Load(p.Context, key)
							return func() (interface{}, error) {
								return thunk()
							}, nil
						},
						// Resolve: resolver.VendorResolver,
					},
				},
			},
		),
		Context: &ctx,
	}
	return &root
}
