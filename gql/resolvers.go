package gql

import (
	"go-graphql-cloud-api/postgres"
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
