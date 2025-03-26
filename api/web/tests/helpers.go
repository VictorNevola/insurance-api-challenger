package tests_test

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type (
	MongoContainer struct {
		*mongodb.MongoDBContainer
		ConnectionString string
	}
)

func createMongodBContainer(ctx context.Context) (*MongoContainer, error) {
	mongoDBContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		panic("Could not start mongo container: " + err.Error())
	}

	connectionString, err := mongoDBContainer.ConnectionString(ctx)
	if err != nil {
		panic("Could not get mongo connection string: " + err.Error())
	}

	return &MongoContainer{
		MongoDBContainer: mongoDBContainer,
		ConnectionString: connectionString,
	}, nil
}

func ConnectionToDB(ctx context.Context, databaseName string) (*mongo.Client, func(), func()) {
	mongoContainer, err := createMongodBContainer(ctx)
	if err != nil {
		panic("error when create mysql container")
	}

	mongoDBClient, err := mongo.Connect(nil, options.Client().ApplyURI(mongoContainer.ConnectionString))
	if err != nil {
		panic(err)
	}

	cleanUp := func() {
		mongoDBClient.Disconnect(ctx)
		mongoContainer.Terminate(ctx)
	}

	clearDataBase := func() {
		clearAllDataBase(ctx, mongoDBClient, databaseName)
	}

	return mongoDBClient, cleanUp, clearDataBase
}

func clearAllDataBase(ctx context.Context, db *mongo.Client, databaseName string) {
	fmt.Println("Cleaning database...")

	collections := []string{"partners", "quotes", "policies"}

	for _, collection := range collections {
		_, err := db.Database(databaseName).Collection(collection).DeleteMany(ctx, bson.M{})
		if err != nil {
			fmt.Printf("Error cleaning collection '%s': %v\n", collection, err)
		} else {
			fmt.Printf("Collection '%s' cleaned successfully.\n", collection)
		}
	}
}
