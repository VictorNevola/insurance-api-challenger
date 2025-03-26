package insurance_test

import (
	"errors"
	"main-api/internal/domain/partners"
	"main-api/internal/infra/cache"
	"main-api/internal/infra/http/insurance"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

var (
	baseURL = "http://localhost:8080"
	apiKey  = "secret"
)

func TestAuthenticate(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	defer gock.Off()

	cacheStorage := cache.NewMockCacheStore(ctrl)
	insuranceProviderClient := insurance.NewInsuranceProviderClient(cacheStorage, baseURL, apiKey)

	gock.InterceptClient(insuranceProviderClient.Client)

	cacheStorage.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	t.Run("Should return success jwt token ", func(t *testing.T) {
		defer gock.Clean()

		mock := gock.New(baseURL).
			Post("/auth").
			Reply(200).
			JSON(map[string]interface{}{
				"access_token": "fake-token",
			})

		token, err := insuranceProviderClient.Authenticate()

		assert.NoError(t, err)
		assert.Equal(t, "fake-token", token.AcessToken)
		assert.True(t, mock.Done())
	})

	t.Run("Should open circuit breaker when failed to authenticate", func(t *testing.T) {
		defer gock.Clean()

		mock := gock.New(baseURL).
			Post("/auth").
			ReplyError(assert.AnError)

		_, err := insuranceProviderClient.Authenticate()

		assert.Error(t, err)
		assert.True(t, mock.Done())
	})
}

func TestCreateQuotation(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	defer gock.Off()

	cacheStorage := cache.NewMockCacheStore(ctrl)
	insuranceProviderClient := insurance.NewInsuranceProviderClient(cacheStorage, baseURL, apiKey)

	gock.InterceptClient(insuranceProviderClient.Client)
	gock.New(baseURL).
		Post("/auth").
		Reply(200).
		JSON(map[string]interface{}{
			"access_token": "fake-token",
		})

	cacheStorage.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	cacheStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	t.Run("Should return success when create new quotation", func(t *testing.T) {
		defer gock.Clean()

		randomUUID := uuid.New()

		gock.New(baseURL).
			Post("/quotation").
			Reply(200).
			JSON(map[string]interface{}{
				"id":        randomUUID,
				"age":       20,
				"price":     120.99,
				"sex":       "F",
				"expire_at": "2025-03-21",
			})

		response, err := insuranceProviderClient.CreateQuotation(
			t.Context(),
			partners.InsuranceProviderCreateQuotationRequest{
				Age: 20,
				Sex: "f/F",
			})

		expectedResponse := &partners.InsuranceProviderCreateQuotationResponse{
			ProviderID: randomUUID,
			Age:        20,
			Price:      120.99,
			Sex:        "F",
			ExpiresAt:  "2025-03-21",
		}

		assert.NoError(t, err)
		assert.EqualValues(t, expectedResponse, response)
	})

	t.Run("Should return error when failed to create quotation", func(t *testing.T) {
		defer gock.Clean()

		gock.New(baseURL).
			Post("/quotation").
			ReplyError(assert.AnError)

		response, err := insuranceProviderClient.CreateQuotation(
			t.Context(),
			partners.InsuranceProviderCreateQuotationRequest{
				Age: 20,
			},
		)

		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestCreatPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	defer gock.Off()

	cacheStorage := cache.NewMockCacheStore(ctrl)
	insuranceProviderClient := insurance.NewInsuranceProviderClient(cacheStorage, baseURL, apiKey)

	gock.InterceptClient(insuranceProviderClient.Client)
	gock.New(baseURL).
		Post("/auth").
		Reply(200).
		JSON(map[string]interface{}{
			"access_token": "fake-token",
		})

	cacheStorage.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	cacheStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	t.Run("Should return success when create new policy", func(t *testing.T) {
		defer gock.Clean()

		gock.New(baseURL).
			Post("/policies").
			Reply(200).
			JSON(map[string]interface{}{
				"id":            "8ed8fe32-a8ce-46c1-aef4-64019eb5a859",
				"quotation_id":  "8ed8fe32-a8ce-46c1-aef4-64019eb5a858",
				"name":          "quotation-test",
				"sex":           "F",
				"date_of_birth": "1998-09-28",
			})

		response, err := insuranceProviderClient.CreatePolicy(t.Context(), partners.InsuranceProviderCreatePolicyRequest{
			QuotationID: uuid.MustParse("8ed8fe32-a8ce-46c1-aef4-64019eb5a858"),
			Name:        "quotation-test",
			Sex:         "F",
			DateOfBirth: "1998-09-28",
		})

		expectedResponse := &partners.InsuranceProviderCreatePolicyResponse{
			ID:          uuid.MustParse("8ed8fe32-a8ce-46c1-aef4-64019eb5a859"),
			QuotationID: uuid.MustParse("8ed8fe32-a8ce-46c1-aef4-64019eb5a858"),
			Name:        "quotation-test",
			Sex:         "F",
			DateOfBirth: "1998-09-28",
		}

		assert.NoError(t, err)
		assert.EqualValues(t, expectedResponse, response)
	})

	t.Run("Should return error when failed to create policies", func(t *testing.T) {
		defer gock.Clean()

		type ErrorScenarios struct {
			ErrorMessage string
			StatusCode   int
		}

		possibleErrors := []ErrorScenarios{
			{
				StatusCode:   400,
				ErrorMessage: "The field 'sex' doesn't match with quotations'",
			},
			{
				StatusCode:   400,
				ErrorMessage: "The field 'date_of_birth' doesn't match with quotations' age",
			},
			{
				StatusCode:   400,
				ErrorMessage: "The quotation was expired",
			},
			{
				StatusCode:   404,
				ErrorMessage: "quotation not found",
			},
		}

		for _, scenario := range possibleErrors {
			gock.New(baseURL).
				Post("/policies").
				Reply(scenario.StatusCode).
				SetError(errors.New(scenario.ErrorMessage))

			response, err := insuranceProviderClient.CreatePolicy(t.Context(), partners.InsuranceProviderCreatePolicyRequest{
				QuotationID: uuid.MustParse("8ed8fe32-a8ce-46c1-aef4-64019eb5a858"),
				Name:        "quotation-test",
				Sex:         "f/F",
				DateOfBirth: "1998-09-28",
			})

			assert.Error(t, errors.New(scenario.ErrorMessage), err)
			assert.Nil(t, response)
		}
	})
}

func TestGetPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	defer gock.Off()

	cacheStorage := cache.NewMockCacheStore(ctrl)
	insuranceProviderClient := insurance.NewInsuranceProviderClient(cacheStorage, baseURL, apiKey)

	gock.InterceptClient(insuranceProviderClient.Client)
	gock.New(baseURL).
		Post("/auth").
		Reply(200).
		JSON(map[string]interface{}{
			"access_token": "fake-token",
		})

	cacheStorage.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	cacheStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	t.Run("Should return success when get policy", func(t *testing.T) {
		defer gock.Clean()

		policyID := "8ed8fe32-a8ce-46c1-aef4-64019eb5a859"

		gock.New(baseURL).
			Get("/policies/" + policyID).
			Reply(200).
			JSON(map[string]interface{}{
				"id":            policyID,
				"quotation_id":  "8ed8fe32-a8ce-46c1-aef4-64019eb5a858",
				"name":          "quotation-test",
				"sex":           "f/F",
				"date_of_birth": "1998-09-28",
			})

		response, err := insuranceProviderClient.GetPolicy(t.Context(), policyID)

		assert.NoError(t, err)
		assert.EqualValues(t, &partners.InsuranceProviderCreatePolicyResponse{
			ID:          response.ID,
			QuotationID: response.QuotationID,
			Name:        response.Name,
			Sex:         response.Sex,
			DateOfBirth: response.DateOfBirth,
		}, response)
	})

	t.Run("Should return not found error when not exists policy", func(t *testing.T) {
		defer gock.Clean()

		policyID := "8ed8fe32-a8ce-46c1-aef4-64019eb5a859"

		gock.New(baseURL).
			Get("/policies/" + policyID).
			Reply(http.StatusNotFound).
			SetError(errors.New("policy not found"))

		response, err := insuranceProviderClient.GetPolicy(t.Context(), policyID)

		assert.Nil(t, response)
		assert.Error(t, errors.New("policy not found"), err)
	})
}
