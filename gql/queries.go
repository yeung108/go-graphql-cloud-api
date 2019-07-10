package gql

import (
	"context"

	"go-graphql-cloud-api/postgres"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
)

// Root holds a pointer to a graphql object
type Root struct {
	Query    *graphql.Object
	Mutation *graphql.Object
	Context  *context.Context
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
		Mutation: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "Mutation",
				Fields: graphql.Fields{
					"editVendor": &graphql.Field{
						// Slice of Vendor type which can be found in types.go
						Type: Vendor,
						Args: graphql.FieldConfigArgument{
							"id": &graphql.ArgumentConfig{
								Type: graphql.String,
							},
						},
						Resolve: resolver.EditVendorResolver,
					},
				},
			},
		),
		Context: &ctx,
	}
	return &root
}
