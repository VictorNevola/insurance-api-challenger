package partners_test

import (
	"context"
	partnersHandler "main-api/api/web/partners"
	tests_test "main-api/api/web/tests"
	partnersDomain "main-api/internal/domain/partners"
	mocks "main-api/internal/infra/repository/mocks"
	partnersRepo "main-api/internal/infra/repository/partners"
	policiesRepo "main-api/internal/infra/repository/policies"
	quotesRepo "main-api/internal/infra/repository/quotes"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	testHelpers struct {
		ctx                     *context.Context
		DBclient                *mongo.Client
		InsuranceProviderClient *mocks.MockInsuranceProvider
	}
)

var (
	helpers      *testHelpers
	databaseName = "test-DB"
)

func testContext(ctrlGoMock *gomock.Controller) (context.Context, *fiber.App, func(), func()) {
	ctx := context.TODO()
	app := fiber.New()

	mongoDBConnection, closeDbConnection, clearAllDataBase := tests_test.ConnectionToDB(
		ctx,
		databaseName,
	)

	partnersRepository := partnersRepo.NewRepo(mongoDBConnection, databaseName)
	quotesRepository := quotesRepo.NewRepo(mongoDBConnection, databaseName)
	policiesRepository := policiesRepo.NewRepo(mongoDBConnection, databaseName)
	insuranceProviderClient := mocks.NewMockInsuranceProvider(ctrlGoMock)

	partnersService := partnersDomain.NewService(partnersDomain.ServiceParams{
		PartnerRepo:             partnersRepository,
		QuoteRepo:               quotesRepository,
		PolicyRepo:              policiesRepository,
		InsuranceClientProvider: insuranceProviderClient,
	})

	partnersHandler.NewHTTPHandler(app, partnersService)

	clearEnviroment := func() {
		closeDbConnection()
		app.Shutdown()
	}

	helpers = &testHelpers{
		DBclient:                mongoDBConnection,
		ctx:                     &ctx,
		InsuranceProviderClient: insuranceProviderClient,
	}

	return ctx, app, clearEnviroment, clearAllDataBase
}

func createAFakePartner() partnersDomain.PartnerEntity {
	entity := partnersDomain.PartnerEntity{
		Name:      "test-patner",
		Cnpj:      "69766865000160",
		CreatedAt: time.Now(),
	}

	result, err := helpers.DBclient.Database(databaseName).
		Collection(partnersRepo.CollectionName).
		InsertOne(*helpers.ctx, entity)
	if err != nil {
		panic("failed to create partner")
	}

	objectID, _ := result.InsertedID.(bson.ObjectID)

	entity.SetID(objectID.Hex())

	return entity
}

func createAFakePolicy(partnerID string) partnersDomain.PolicyEntity {
	entity := partnersDomain.PolicyEntity{
		QuotationID: uuid.New(),
		ProviderID:  uuid.New(),
		Sex:         "M",
		Name:        "test-policy",
		DateOfBirth: "1998-09-28",
		PartnerID:   partnerID,
	}

	result, err := helpers.DBclient.Database(databaseName).
		Collection(policiesRepo.CollectionName).
		InsertOne(*helpers.ctx, map[string]interface{}{
			"provider_id":   entity.ProviderID.String(),
			"quotation_id":  entity.QuotationID.String(),
			"partner_id":    entity.PartnerID,
			"name":          entity.Name,
			"sex":           entity.Sex,
			"date_of_birth": entity.DateOfBirth,
		})
	if err != nil {
		panic("failed to create partner")
	}

	objectID, _ := result.InsertedID.(bson.ObjectID)

	entity.ID = objectID.Hex()

	return entity

}

func setResponseInsuranceQuotation(
	dataReturn partnersDomain.InsuranceProviderCreateQuotationResponse,
) {
	helpers.InsuranceProviderClient.EXPECT().
		CreateQuotation(gomock.Any(), gomock.Any()).
		Return(&dataReturn, nil)
}

func setResponseInsurancePolicy(
	dataReturn partnersDomain.InsuranceProviderCreatePolicyResponse,
) {
	helpers.InsuranceProviderClient.EXPECT().
		CreatePolicy(gomock.Any(), gomock.Any()).
		Return(&dataReturn, nil)
}

func setResponseInsurancePolicyError() {
	helpers.InsuranceProviderClient.EXPECT().
		CreatePolicy(gomock.Any(), gomock.Any()).
		Return(nil, fiber.NewError(fiber.StatusBadRequest))
}

func setResponseGetInsurancePolicy(
	dataReturn partnersDomain.InsuranceProviderCreatePolicyResponse,
) {
	helpers.InsuranceProviderClient.EXPECT().
		GetPolicy(gomock.Any(), gomock.Any()).
		Return(&dataReturn, nil)
}
