package quotes

import (
	"context"
	"fmt"
	"main-api/internal/domain/partners"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	Repo struct {
		DatabaseName string
		DB           *mongo.Client
	}
)

var (
	collectionName = "quotes"
)

func NewRepo(db *mongo.Client, dbName string) *Repo {
	return &Repo{
		DatabaseName: dbName,
		DB:           db,
	}
}

func (r *Repo) Create(ctx context.Context, quote *partners.QuoteEntity) error {
	collection := r.DB.Database(r.DatabaseName).Collection(collectionName)

	result, err := collection.InsertOne(ctx, map[string]interface{}{
		"provider_id": quote.ProviderID.String(),
		"partner_id":  quote.PartnerID,
		"age":         quote.Age,
		"sex":         quote.Sex,
		"price":       quote.Price,
		"expires_at":  quote.ExpiresAt,
		"created_at":  quote.CreatedAt,
	})
	if err != nil {
		return err
	}

	objectID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return fmt.Errorf("error on convert inserted id to ObjectID")
	}

	quote.ID = objectID.Hex()

	return nil
}
