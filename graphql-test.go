package main

import (
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type query struct{}
type user struct{}
type request struct{}

func (_ *query) Hello() string { return "Hello, world!" }
func (_ *query) Bye() string   { return "Bye, world!" }

func (_ *user) Name() string        { return "John" }
func (_ *user) Email() string       { return "ABC" }
func (_ *user) PhoneNumber() string { return "DEF" }
func (_ *user) Address() string     { return "GHI" }

func (_ *request) Name() string         { return "Albert Gator" }
func (_ *request) StoreName() string    { return "Publix" }
func (_ *request) StoreAddress() string { return "Gainesville, FL" }
func (_ *request) Items() string        { return "Pub Sub" }

func main() {
	u := `
			type Query {
					name: String!
					email: String!
					phoneNumber: String!
					address: String!
			}
		`
	r := `
			type Query {
				name: String!
				storeName: String!
				storeAddress: String!
				items: String!
			}
	`
	userSchema := graphql.MustParseSchema(u, &user{})
	requestSchema := graphql.MustParseSchema(r, &request{})
	http.Handle("/user", &relay.Handler{Schema: userSchema})
	http.Handle("/request", &relay.Handler{Schema: requestSchema})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
