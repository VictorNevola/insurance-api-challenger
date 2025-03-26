package partners

import (
	"context"
	"errors"
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
	CollectionName = "partners"
)

func NewRepo(db *mongo.Client, dbName string) *Repo {
	return &Repo{
		DatabaseName: dbName,
		DB:           db,
	}
}

func (r *Repo) GetByFilter(ctx context.Context, filter map[string]interface{}) (*partners.PartnerEntity, error) {
	collection := r.DB.Database(r.DatabaseName).Collection(CollectionName)
	partner := new(partners.PartnerEntity)

	err := collection.FindOne(ctx, filter).Decode(partner)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return partner, nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*partners.PartnerEntity, error) {
	collection := r.DB.Database(r.DatabaseName).Collection(CollectionName)
	partner := new(partners.PartnerEntity)

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, map[string]interface{}{"_id": objectID}).Decode(partner)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	partner.SetID(id)

	return partner, nil
}

func (r *Repo) Create(ctx context.Context, partner *partners.PartnerEntity) error {
	collection := r.DB.Database(r.DatabaseName).Collection(CollectionName)

	result, err := collection.InsertOne(ctx, map[string]interface{}{
		"name":       partner.Name,
		"cnpj":       partner.Cnpj,
		"created_at": partner.CreatedAt,
	})
	if err != nil {
		return err
	}

	objectID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return fmt.Errorf("error on convert inserted id to ObjectID")
	}

	partner.SetID(objectID.Hex())

	return nil
}
