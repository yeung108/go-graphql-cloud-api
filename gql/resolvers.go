package gql

import (
	"go-graphql-cloud-api/postgres"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// Resolver struct holds a connection to our database
type Resolver struct {
	db *postgres.Db
}

// VendorResolver resolves our settings query through a db call to GetSettings
// func (r *Resolver) VendorResolver(p graphql.ResolveParams) (interface{}, error) {
// 	vendors := r.db.GetVendorProducts()
// 	return vendors, nil
// }
func (r *Resolver) VendorResolver(p graphql.ResolveParams) (interface{}, error) {
	var (
		v       = p.Context.Value
		c       = v("client").(*Client)
		loaders = v("loaders").(map[string]*dataloader.Loader)
		id      = p.Args["id"].(string)
		key     = NewResolverKey(id, c)
	)
	thunk := loaders["GetVendors"].Load(p.Context, key)
	return func() (interface{}, error) {
		return thunk()
	}, nil
}

func (r *Resolver) EditVendorResolver(p graphql.ResolveParams) (interface{}, error) {
	k, _ := uuid.FromString(p.Args["id"].(string))
	err := r.db.EditVendors(k)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
