package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "github.com/graphql-go/graphql"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
	Email string `bson:"email,omitempty"`
	PhoneNumber string `bson:"phoneNumber,omitempty"`
	Address string `bson:"address,omitempty"`
}

type Request struct {
	user User
	storeAddress string
	items string
}

type Store struct {
	name string
	address string
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://15dani1:hacknow@cluster0-f47on.gcp.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	hacknowDatabase := client.Database("hacknow")
	usersCollection := hacknowDatabase.Collection("users")
	requestsCollection := hacknowDatabase.Collection("requests")
	rohil := User{
		Name: "Rohil Tuli", 
		Email: "rohil.tuli@gmail.com", 
		PhoneNumber: "8133732574", 
		Address: "4100 George J Bean Pkwy, Tampa, FL 33607",
	}
	userResult, err := usersCollection.InsertOne(ctx, rohil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(userResult.InsertedID)

	requestResult, err := requestsCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{"user", rohil},
			{"items", "testing"},
		},
		bson.D{
			{"user", rohil},
			{"items", "testing again"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(requestResult.InsertedIDs)
}
