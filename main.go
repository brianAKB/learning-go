package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	DatabaseName   = "sample-db"
	DatabaseURL    = "mongodb://localhost:27017"
	CollectionName = "sample-collection"
)

type SampleNestedData struct {
	InnerStrValue string `bson:"innerStrValue"`
}

type SampleData struct {
	ID          uint64           `bson:"id,omitempty"`
	StrValue    string           `bson:"strValue"`
	IntValue    int              `bson:"intValue"`
	BoolValue   bool             `bson:"boolValue"`
	NestedValue SampleNestedData `bson:"nestedValue"`
}

// newDatabaseClient establishes and verifies a connection with the target database.
func newDatabaseClient(databaseURL string, timeout int) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(databaseURL)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func InsertItem(col *mongo.Collection, data SampleData) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := col.InsertOne(ctx, &data)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func GetItem(col *mongo.Collection, filter bson.D) SampleData {
	var data SampleData

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	found := col.FindOne(ctx, filter)
	if err := found.Decode(&data); err != nil {
		return SampleData{}
	}

	return data
}

func UpdateItem(col *mongo.Collection, filter bson.D, updateParams bson.D) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := col.UpdateOne(ctx, filter, updateParams, options.Update())
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func DeleteItem(col *mongo.Collection, filter bson.D) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func main() {

	// create a new client
	client, _ := newDatabaseClient(DatabaseURL, 10)

	// retrieve the DB and collection from the client
	db := client.Database(DatabaseName)
	collection := db.Collection(CollectionName)

	/*
		data := SampleData{
			ID:        1,
			StrValue:  "a string in the doc",
			IntValue:  42,
			BoolValue: true,
			NestedValue: SampleNestedData{
				InnerStrValue: "a string in a nested doc",
			},
		}
	*/
	//InsertItem(collection, data)

	searchFilter := bson.D{
		{Key: "id", Value: 1},
	}

	retrievedData := GetItem(collection, searchFilter)

	updateParameters := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "boolvalue", Value: false},
		}},
	}

	UpdateItem(collection, searchFilter, updateParameters)

	DeleteItem(collection, searchFilter)

	fmt.Printf("recovered doc: %v", retrievedData)
	element := bson.E{Key: "id", Value: "SampleID"}

	filter := bson.D{
		element,
	}
	fmt.Printf("%v", filter)

}
