package main

import (
	"context"
	"fmt"
	"strings"
	"log"
	"net/http"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

// type query struct{}
// type user struct{

// }

func getCoordinates(address string) string {
	resp, err := http.Get("https://api.mapbox.com/geocoding/v5/mapbox.places/" + address + ".json?access_token=pk.eyJ1IjoiMTVkYW5pMSIsImEiOiJjazlmNWdvdG4wMGVvM2xubjdqcTducXM1In0.H7cu4oj3nkFtR23KeFEliQ")

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string][](map[string](map[string][]float64))
	err = json.Unmarshal(body, &result)
	// var result2 map[string]interface{}
	// err = json.Unmarshal(result, &result2)
	// log.Println(string(body))
	log.Println(result["features"][0]["geometry"]["coordinates"])
	lat := strconv.FormatFloat(result["features"][0]["geometry"]["coordinates"][0], 'f', -1, 64)
	//lat := (result["features"][0]["geometry"]["coordinates"][0]).string
	long := strconv.FormatFloat(result["features"][0]["geometry"]["coordinates"][1], 'f', -1, 64)
	//long := (result["features"][0]["geometry"]["coordinates"][1]).string
	return lat + "," + long
}

func mapBoxDriver(address1, address2 string) string {

	address1New := strings.ReplaceAll(address1, " ", "%20")
	address2New := strings.ReplaceAll(address2, " ", "%20")
	requestBody := "coordinates=" + getCoordinates(address1New) + ";" + getCoordinates(address2New)

	resp, err := http.Post("https://api.mapbox.com/directions/v5/mapbox/driving?access_token=pk.eyJ1IjoiMTVkYW5pMSIsImEiOiJjazlmNWdvdG4wMGVvM2xubjdqcTducXM1In0.H7cu4oj3nkFtR23KeFEliQ", "application/x-www-form-urlencoded", bytes.NewBufferString(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
	return string(body)
	// log.Println(resp)
}


type Account struct {
	NameField        string `bson:"name,omitempty"`
	AddressField     string `bson:"address,omitempty"`
	PhoneNumberField string `bson:"phoneNumber,omitempty"`
}

type Request struct {
	ItemsField   string   `bson:"items,omitempty"`
	AccountField *Account `bson:"account,omitempty"`
}

// type request struct{}



type AccountResolver struct{}
type RequestResolver struct{}

func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // allow cross domain AJAX requests
        w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		w.WriteHeader(http.StatusOK)
        next.ServeHTTP(w,r)
    })
}

// func (_ *query) Hello() string { return "Hello, world!" }
// func (_ *query) Bye() string   { return "Bye, world!" }
func (r AccountResolver) Request(ctx context.Context, args struct{ StoreAddress string }) *Request {
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
	requestsCollection := hacknowDatabase.Collection("requests")
	var result Request
	err = requestsCollection.FindOne(ctx, bson.D{{"store", args.StoreAddress}}).Decode(&result)
	if err != nil {
		return &Request{}
	}
	return &result
// 	return []Request{&Request{}, &Request{}}
}

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

type PhoneStruct struct {
	PhoneNumber string
}

func (r AccountResolver) CreateRequest(ctx context.Context, args struct {
	StoreAddress string
	PhoneNumber  string
	Items        string
}) *Request {
	return &Request{ItemsField: args.Items, AccountField: r.Account(ctx, PhoneStruct{args.PhoneNumber})}
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

func (a *Request) Items() string {
	return a.ItemsField
}

func (a *Request) Account() *Account {
	return a.AccountField
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
			createRequest(StoreAddress: String!, PhoneNumber: String!, Items: String!): Request
		}
		type Query {
			account(PhoneNumber: String!): Account
			request(StoreAddress: String!): Request
		}
		type Account {
			name: String!
			address: String!
			phoneNumber: String!
		}
		type Request {
			items: String!
			account: Account!
		}
		`
	// request(PhoneNumber: String!): [Request]
	// r := `
	// 	schema {
	// 		query: Query
	// 		mutation: Mutation
	// 	}
	// 	type Mutation {
	// 		createRequest(StoreAddress: String!, PhoneNumber: String!): Request
	// 	}
	// 	type Query {
	// 		request(PhoneNumber: String!): [Request]
	// 	}
	// 	type Request {
	// 		items: String!
	// 		account: Account!
	// 	}
	// 	type Account {
	// 		name: String!
	// 		address: String!
	// 		phoneNumber: String!
	// 	}
	// `
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
	log.Println(mapBoxDriver("Land O Lakes", "Orlando"))
	accountSchema := graphql.MustParseSchema(a, &AccountResolver{})
	//requestSchema := graphql.MustParseSchema(r, &AccountResolver{})
	http.Handle("/account", CorsMiddleware(&relay.Handler{Schema: accountSchema}))
	//http.Handle("/request", &relay.Handler{Schema: requestSchema})
	log.Fatal(http.ListenAndServe(":5000", nil))
}
