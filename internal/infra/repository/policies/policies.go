package policies

import (
	"context"
	"fmt"
	"main-api/internal/domain/partners"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	Repo struct {
		DatabaseName string
		DB           *mongo.Client
	}

	policyResultDB struct {
		ID          bson.ObjectID `bson:"_id"`
		ProviderID  string        `bson:"provider_id"`
		QuotationID string        `bson:"quotation_id"`
		PartnerID   string        `bson:"partner_id"`
		Name        string        `bson:"name"`
		Sex         string        `bson:"sex"`
		DateOfBirth string        `bson:"date_of_birth"`
	}
)

var (
	CollectionName = "policies"
)

func NewRepo(db *mongo.Client, dbName string) *Repo {
	return &Repo{
		DatabaseName: dbName,
		DB:           db,
	}
}

func (r *Repo) Create(ctx context.Context, policy *partners.PolicyEntity) error {
	collection := r.DB.Database(r.DatabaseName).Collection(CollectionName)

	result, err := collection.InsertOne(ctx, map[string]interface{}{
		"provider_id":   policy.ProviderID.String(),
		"quotation_id":  policy.QuotationID.String(),
		"partner_id":    policy.PartnerID,
		"name":          policy.Name,
		"sex":           policy.Sex,
		"date_of_birth": policy.DateOfBirth,
	})
	if err != nil {
		return err
	}

	objectID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return fmt.Errorf("error on convert inserted id to ObjectID")
	}

	policy.ID = objectID.Hex()

	return nil
}

func (r *Repo) GetByIdAndPartnerID(ctx context.Context, policyID, partnerID string) (*partners.PolicyEntity, error) {
	collection := r.DB.Database(r.DatabaseName).Collection(CollectionName)
	id, err := bson.ObjectIDFromHex(policyID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        id,
		"partner_id": partnerID,
	}

	var result policyResultDB
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, partners.ErrPolicyNotFound
		}
		return nil, err
	}

	return &partners.PolicyEntity{
		ID:          result.ID.Hex(),
		PartnerID:   result.PartnerID,
		Name:        result.Name,
		DateOfBirth: result.DateOfBirth,
		QuotationID: uuid.MustParse(result.QuotationID),
		ProviderID:  uuid.MustParse(result.ProviderID),
		Sex:         partners.SexEnum(result.Sex),
	}, nil
}
