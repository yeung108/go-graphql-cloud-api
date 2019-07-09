package gql

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

// ExecuteQuery runs our graphql queries
func ExecuteQuery(query string, schema graphql.Schema, ctx context.Context) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
	})

	b, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[GraphQL] result: %s", b)

	// Error check
	if len(result.Errors) > 0 {
		fmt.Printf("Unexpected errors inside ExecuteQuery: %v", result.Errors)
	}

	return result
}
