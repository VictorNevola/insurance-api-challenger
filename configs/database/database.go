package database

import (
	"main-api/configs/envs"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitMongoDB() *mongo.Client {
	mongoDBClient, err := mongo.Connect(nil, options.Client().ApplyURI(envs.AppConfig.MongoURL))
	if err != nil {
		panic(err)
	}

	return mongoDBClient
}
