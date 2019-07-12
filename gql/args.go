package gql

import (
	"github.com/graphql-go/graphql"
)

var LanguageJsonArgs = graphql.InputObjectConfig{
	Fields: graphql.InputObjectConfigFieldMap{
		"en": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"zh": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
}

var VendorArgs = graphql.InputObjectConfig{
	Name: "VendorArgs",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"description": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"products": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(ProductArgs),
		},
		"stores": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(StoreArgs),
		},
	},
}

// Product describes a graphql args containing a Product
var ProductArgs = graphql.InputObjectConfig{
	Name: "ProductArgs",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"photo": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"code": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"is_virtual_product": &graphql.InputObjectFieldConfig{
			Type: graphql.Boolean,
		},
		"barcode": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"descriptions": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(LanguageJsonArgs),
		},
		"brand_names": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(LanguageJsonArgs),
		},
		"names": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(LanguageJsonArgs),
		},
		"optional_data": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(LanguageJsonArgs),
		},
	},
}

// Store describes a graphql args containing a Store
var StoreArgs = graphql.InputObjectConfig{
	Name: "StoreArgs",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"created_at": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"updated_at": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"mongo_id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"code": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"model": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"address": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"last_online_at": &graphql.InputObjectFieldConfig{
			Type: graphql.DateTime,
		},
		"last_get": &graphql.InputObjectFieldConfig{
			Type: graphql.DateTime,
		},
		"last_refill": &graphql.InputObjectFieldConfig{
			Type: graphql.DateTime,
		},
		"last_reset": &graphql.InputObjectFieldConfig{
			Type: graphql.DateTime,
		},
		"last_sync": &graphql.InputObjectFieldConfig{
			Type: graphql.DateTime,
		},
		"unsubmitted_order_count": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
	},
}
