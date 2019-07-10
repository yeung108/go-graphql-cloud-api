package gql

import (
	"go-graphql-cloud-api/postgres"
	"go-graphql-cloud-api/scalar"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
)

type Client struct {
	Resolver *Resolver
}

type ResolverKey struct {
	Key    string
	Client *Client
}

func (rk *ResolverKey) client() *Client {
	return rk.Client
}

func (client *Client) resolver() *Resolver {
	return client.Resolver
}

func NewResolverKey(key string, client *Client) *ResolverKey {
	return &ResolverKey{
		Key:    key,
		Client: client,
	}
}

func (rk *ResolverKey) String() string {
	return rk.Key
}

func (rk *ResolverKey) Raw() interface{} {
	return rk.Key
}

var LanguageJson = graphql.NewObject(graphql.ObjectConfig{
	Name: "LanguageJson",
	Fields: graphql.Fields{
		"en": &graphql.Field{Type: scalar.NullScalar},
		"zh": &graphql.Field{Type: scalar.NullScalar},
	},
})

var Vendor = graphql.NewObject(graphql.ObjectConfig{
	Name: "Vendor",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: scalar.NullScalar,
		},
		"created_at": &graphql.Field{
			Type: scalar.NullScalar,
		},
		"updated_at": &graphql.Field{
			Type: scalar.NullScalar,
		},
		"mongo_id": &graphql.Field{
			Type: scalar.NullScalar,
		},
		"name": &graphql.Field{
			Type: scalar.NullScalar,
		},
		"description": &graphql.Field{
			Type: scalar.NullScalar,
		},
		"products": &graphql.Field{
			Type: graphql.NewList(Product),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var (
					v       = p.Context.Value
					c       = v("client").(*Client)
					loaders = v("loaders").(map[string]*dataloader.Loader)
					vendor  = p.Source.(postgres.Vendor)
					key     = NewResolverKey(vendor.ID.String(), c)
				)
				thunk := loaders["GetVendorProducts"].Load(p.Context, key)
				return func() (interface{}, error) {
					return thunk()
				}, nil
			},
		},
		"stores": &graphql.Field{
			Type: graphql.NewList(Store),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var (
					v       = p.Context.Value
					c       = v("client").(*Client)
					loaders = v("loaders").(map[string]*dataloader.Loader)
					vendor  = p.Source.(postgres.Vendor)
					key     = NewResolverKey(vendor.ID.String(), c)
				)
				thunk := loaders["GetVendorStores"].Load(p.Context, key)
				return func() (interface{}, error) {
					return thunk()
				}, nil
			},
		},
	},
})

// Product describes a graphql object containing a Product
var Product = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"created_at": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"updated_at": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"mongo_id": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"photo": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"code": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"is_virtual_product": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"barcode": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"descriptions": &graphql.Field{
				Type: LanguageJson,
			},
			"brand_names": &graphql.Field{
				Type: LanguageJson,
			},
			"names": &graphql.Field{
				Type: LanguageJson,
			},
			"optional_data": &graphql.Field{
				Type: LanguageJson,
			},
			"vendor_id": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"supplier_id": &graphql.Field{
				Type: scalar.NullScalar,
			},
		},
	},
)

// Store describes a graphql object containing a Store
var Store = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Store",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"created_at": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"updated_at": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"mongo_id": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"code": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"name": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"model": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"address": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"last_online_at": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"last_get": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"last_refill": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"last_reset": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"last_sync": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"unsubmitted_order_count": &graphql.Field{
				Type: scalar.NullScalar,
			},
			"vendor_id": &graphql.Field{
				Type: scalar.NullScalar,
			},
		},
	},
)
