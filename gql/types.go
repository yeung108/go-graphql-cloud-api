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
		"en": &graphql.Field{Type: graphql.String},
		"zh": &graphql.Field{Type: graphql.String},
	},
})

var Vendor = graphql.NewObject(graphql.ObjectConfig{
	Name: "Vendor",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updated_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"mongo_id": &graphql.Field{
			Type: scalar.NullStringScalar,
		},
		"name": &graphql.Field{
			Type: scalar.NullStringScalar,
		},
		"description": &graphql.Field{
			Type: scalar.NullStringScalar,
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
	},
})

// Product describes a graphql object containing a Product
var Product = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"mongo_id": &graphql.Field{
				Type: scalar.NullStringScalar,
			},
			"photo": &graphql.Field{
				Type: scalar.NullStringScalar,
			},
			"code": &graphql.Field{
				Type: scalar.NullStringScalar,
			},
			"is_virtual_product": &graphql.Field{
				Type: graphql.Boolean,
			},
			"barcode": &graphql.Field{
				Type: scalar.NullStringScalar,
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
				Type: scalar.NullStringScalar,
			},
			"supplier_id": &graphql.Field{
				Type: scalar.NullStringScalar,
			},
		},
	},
)

// Store describes a graphql object containing a Store
var Store = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: scalar.SpecialDateScalar,
			},
			"updated_at": &graphql.Field{
				Type: scalar.SpecialDateScalar,
			},
			"mongo_id": &graphql.Field{
				Type: scalar.NullStringScalar,
			},
			"code": &graphql.Field{
				Type: scalar.NullStringScalar,
			},
			"name": &graphql.Field{
				Type: scalar.NullStringScalar,
			},
			"model": &graphql.Field{
				Type: LanguageJson,
			},
			"address": &graphql.Field{
				Type: LanguageJson,
			},
			"names": &graphql.Field{
				Type: LanguageJson,
			},
			"optional_data": &graphql.Field{
				Type: LanguageJson,
			},
			"vendor_id": &graphql.Field{
				Type: scalar.NullStringScalar,
			},
			"supplier_id": &graphql.Field{
				Type: scalar.NullStringScalar,
			},
		},
	},
)
