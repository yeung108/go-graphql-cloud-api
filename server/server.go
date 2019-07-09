package server

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"go-graphql-cloud-api/gql"
	"net/http"

	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"
)

// Server will hold connection to the db as well as handlers
type Server struct {
	GqlSchema *graphql.Schema
	Context   *context.Context
}

type reqBody struct {
	Query     string `json:"query"`
	Signature string `json:"signature"`
}

func HashSha256(msg string) string {
	h := sha256.New()
	h.Write([]byte(msg))
	bs := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(bs)
}

// GraphQL returns an http.HandlerFunc for our /graphql endpoint
func (s *Server) GraphQL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check to ensure query was provided in the request body
		if r.Body == nil {
			http.Error(w, "Must provide graphql query in request body", 400)
			return
		}

		var rBody reqBody
		// Decode the request body into rBody
		err := json.NewDecoder(r.Body).Decode(&rBody)
		if err != nil {
			http.Error(w, "Error parsing JSON request body", 400)
		}

		// Authentication here
		// auth := rBody.Signature
		// publicKey := ciphers.LoadRSAPublicPemKey("public.pem")
		// verified := ciphers.VerifyWithPublicKey(auth, rBody.Query, *publicKey)

		// fmt.Println("rBody.Query: ", rBody.Query)
		// fmt.Println("Signature: ", auth)
		// fmt.Println("verified: ", verified)

		// Check if signature == query
		// if err != nil || !verified {
		// 	fmt.Println(err)
		// 	http.Error(w, "Authentication Error", 401)
		// } else {
		// Execute graphql query
		result := gql.ExecuteQuery(rBody.Query, *s.GqlSchema, *s.Context)

		// render.JSON comes from the chi/render package and handles
		// marshalling to json, automatically escaping HTML and setting
		// the Content-Type as application/json.
		render.JSON(w, r, result)
		// }
	}
}
