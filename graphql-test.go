package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

// type query struct{}
// type user struct{

// }

type Account struct {
	NameField        string `bson:"name,omitempty"`
	AddressField     string `bson:"address,omitempty"`
	PhoneNumberField string `bson:"phoneNumber,omitempty"`
}

// type request struct{}

type AccountResolver struct{}

// func (_ *query) Hello() string { return "Hello, world!" }
// func (_ *query) Bye() string   { return "Bye, world!" }

func (r AccountResolver) Account(ctx context.Context, args struct{ PhoneNumber string }) *Account {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://15dani1:hacknow@cluster0-f47on.gcp.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	hacknowDatabase := client.Database("hacknow")
	usersCollection := hacknowDatabase.Collection("users")
	var result Account
	err = usersCollection.FindOne(ctx, bson.D{{"phoneNumber", args.PhoneNumber}}).Decode(&result)
	if err != nil {
		return &Account{}
	}
	return &result
}

func (r AccountResolver) CreateAccount(ctx context.Context, args struct {
	Name        string
	Address     string
	PhoneNumber string
}) *Account {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://15dani1:hacknow@cluster0-f47on.gcp.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	hacknowDatabase := client.Database("hacknow")
	usersCollection := hacknowDatabase.Collection("users")
	newAccount := Account{
		NameField:        args.Name,
		AddressField:     args.Address,
		PhoneNumberField: args.PhoneNumber,
	}
	insertionResult, err := usersCollection.InsertOne(ctx, newAccount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertionResult.InsertedID)
	return &newAccount
}

func (a *Account) Name() string {
	return a.NameField
}

func (a *Account) Address() string {
	return a.AddressField
}

func (a *Account) PhoneNumber() string {
	return a.PhoneNumberField
}

// func (_ *user) Name() string        { return "John" }
// func (_ *user) Email() string       { return "ABC" }
// func (_ *user) PhoneNumber() string { return "DEF" }
// func (_ *user) Address() string     { return "GHI" }

// func (_ *request) Name() string         { return "Albert Gator" }
// func (_ *request) StoreName() string    { return "Publix" }
// func (_ *request) StoreAddress() string { return "Gainesville, FL" }
// func (_ *request) Items() string        { return "Pub Sub" }

func main() {
	// u := `
	// 		type Query {
	// 				name: String!
	// 				email: String!
	// 				phoneNumber: String!
	// 				address: String!
	// 		}
	// 	`
	a := `
		schema {
			query: Query
			mutation: Mutation
		}
		type Mutation {
			createAccount(Name: String!, Address: String!, PhoneNumber: String!): Account
		}
		type Query {
			account(PhoneNumber: String!): Account
		}
		type Account {
			name: String!
			address: String!
			phoneNumber: String!
		}
	`
	// r := `
	// 		type Query {
	// 			name: String!
	// 			storeName: String!
	// 			storeAddress: String!
	// 			items: String!
	// 		}`

	// userSchema := graphql.MustParseSchema(u, &user{})
	// requestSchema := graphql.MustParseSchema(r, &request{})
	// http.Handle("/user", &relay.Handler{Schema: userSchema})
	// http.Handle("/request", &relay.Handler{Schema: requestSchema})
	accountSchema := graphql.MustParseSchema(a, &AccountResolver{})
	http.Handle("/account", &relay.Handler{Schema: accountSchema})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
