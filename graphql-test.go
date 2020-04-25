package main

import (
	"log"
	"net/http"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type query struct{}
type user struct{}

func (_ *query) Hello() string { return "Hello, world!" }
func (_ *query) Bye() string { return "Bye, world!" }

func (_ *user) Name() string { return "Byeeee" }
func (_ *user) Email() string { return "ABC" }
func (_ *user) PhoneNumber() string { return "DEF" }
func (_ *user) Address() string { return "GHI" }


func main() {
        s := `
                type Query {
						hello: String!
						bye: String!
                }
		`
		u := `
				type Query {
						name: String!
						email: String!
						phoneNumber: String!
						address: String!
			}
		`
		schema := graphql.MustParseSchema(s, &query{})
        userSchema := graphql.MustParseSchema(u, &user{})
		http.Handle("/query", &relay.Handler{Schema: schema})
        http.Handle("/querytwo", &relay.Handler{Schema: userSchema})
        log.Fatal(http.ListenAndServe(":8080", nil))
}