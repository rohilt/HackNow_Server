package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	quickstartDatabase := client.Database("quickstart")
	podcastsCollection := quickstartDatabase.Collection("podcasts")
	episodeCollection := quickstartDatabase.Collection("episodes")
	podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
		{Key: "title", Value: "The Polyglot Developer Podcast"},
		{Key: "author", Value: "Rahul Dani"},
		//OR you can remove Key and Value and have a comma separated list
		{"tags", bson.A{"development", "programming", "coding"}},
		//bson.A is a bson array
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(podcastResult.InsertedID)

	episodeResult, err := episodeCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{"podcast", podcastResult.InsertedID},
			{"title", "Episode #1"},
			{"description", "This is the first episode."},
			{"duration", 25},
		},
		bson.D{
			{"podcast", podcastResult.InsertedID},
			{"title", "Episode #2"},
			{"description", "This is the second episode."},
			{"duration", 32},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(episodeResult.InsertedIDs)
}
