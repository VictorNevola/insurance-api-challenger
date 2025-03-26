package partners_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	partnersHandler "main-api/api/web/partners"
	partnerDomain "main-api/internal/domain/partners"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	PartnerPath = "/partners/"
)

func TestCreatePartner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, server, cleanUp, clearAllDataBase := testContext(ctrl)
	defer cleanUp()

	t.Run("Should create sucessfuly a partner and return", func(t *testing.T) {
		defer clearAllDataBase()

		payload := map[string]interface{}{
			"name": "180 Seguros",
			"cnpj": "12345678901234",
		}

		jsonData, err := json.Marshal(payload)
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, PartnerPath, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		var response partnersHandler.CreatePartnerResponseData
		err = json.NewDecoder(resp.Body).Decode(&response)

		expectedResponse := partnersHandler.CreatePartnerResponseData{
			ID:        response.ID,
			Name:      "180 Seguros",
			Cnpj:      "12345678901234",
			CreatedAt: response.CreatedAt,
		}

		assert.NoError(t, err)
		assert.EqualValues(t, expectedResponse, response)
	})

	t.Run("Not should create a partner when have an invalid payload", func(t *testing.T) {
		payload := map[string]interface{}{
			"name": "180 Seguros",
		}

		jsonData, err := json.Marshal(payload)
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, PartnerPath, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Not should create a partner with the same CNPJ", func(t *testing.T) {
		defer clearAllDataBase()

		payload := map[string]interface{}{
			"name": "180 Seguros",
			"cnpj": "12345678901234",
		}

		jsonData, err := json.Marshal(payload)
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, PartnerPath, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		req, _ = http.NewRequest(http.MethodPost, PartnerPath, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err = server.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})
}

func TestCreateQuote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, server, cleanUp, clearAllDataBase := testContext(ctrl)
	defer cleanUp()

	t.Run("Should create sucessfuly a quote and return", func(t *testing.T) {
		defer clearAllDataBase()

		fakeInsuranceCreateQuotation := partnerDomain.InsuranceProviderCreateQuotationResponse{
			ProviderID: uuid.New(),
			Age:        10,
			Price:      130.99,
			Sex:        "M",
			ExpiresAt:  "2999-03-24",
		}

		fakePartner := createAFakePartner()
		setResponseInsuranceQuotation(fakeInsuranceCreateQuotation)

		payload := map[string]interface{}{
			"age": 10,
			"sex": "M",
		}

		jsonData, err := json.Marshal(payload)
		assert.NoError(t, err)

		path := fmt.Sprintf("%s%s/quotes", PartnerPath, fakePartner.ID)

		req, _ := http.NewRequest(http.MethodPost, path, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		var response partnersHandler.CreateQuoteResponseData
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.NotEmpty(t, response.ID)
		assert.NotEmpty(t, response.ExpiresAt)
		assert.NotEmpty(t, response.CreatedAt)

		assert.Equal(t, uint(10), response.Age)
		assert.Equal(t, "M", response.Sex)
		assert.Equal(t, fakeInsuranceCreateQuotation.Price, response.Price)
	})

	t.Run("Not should create a quote when have an invalid payload", func(t *testing.T) {
		fakePartner := createAFakePartner()

		payload := map[string]interface{}{
			"age": 9999,
			"sex": "A",
		}

		jsonData, err := json.Marshal(payload)
		assert.NoError(t, err)

		path := fmt.Sprintf("%s%s/quotes", PartnerPath, fakePartner.ID)

		req, _ := http.NewRequest(http.MethodPost, path, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Not should create a quote when patner not exists", func(t *testing.T) {
		defer clearAllDataBase()
		createAFakePartner()

		payload := map[string]interface{}{
			"age": 10,
			"sex": "M",
		}

		jsonData, err := json.Marshal(payload)
		assert.NoError(t, err)

		path := fmt.Sprintf("%s%s/quotes", PartnerPath, "67e1a0349d00cb473900ac08")

		req, _ := http.NewRequest(http.MethodPost, path, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCreatePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, server, cleanUp, clearAllDataBase := testContext(ctrl)
	defer cleanUp()

	fakeInsuranceCreatePolicy := partnerDomain.InsuranceProviderCreatePolicyResponse{
		ID:          uuid.New(),
		QuotationID: uuid.New(),
		Name:        "policy-test",
		Sex:         "F",
		DateOfBirth: "1998-09-28",
	}

	t.Run("Should create sucessfuly a policy and return", func(t *testing.T) {
		defer clearAllDataBase()

		fakePartner := createAFakePartner()
		setResponseInsurancePolicy(fakeInsuranceCreatePolicy)

		payload := map[string]interface{}{
			"quotation_id":  fakeInsuranceCreatePolicy.QuotationID,
			"name":          fakeInsuranceCreatePolicy.Name,
			"sex":           fakeInsuranceCreatePolicy.Sex,
			"date_of_birth": fakeInsuranceCreatePolicy.DateOfBirth,
		}

		jsonData, err := json.Marshal(payload)
		assert.NoError(t, err)

		path := fmt.Sprintf("%s%s/policies", PartnerPath, fakePartner.ID)

		req, _ := http.NewRequest(http.MethodPost, path, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		var response partnersHandler.CreatePolicyResponseData
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, response.ID, fakeInsuranceCreatePolicy.ID.String())
		assert.Equal(t, response.Sex, fakeInsuranceCreatePolicy.Sex)
		assert.Equal(t, response.Name, fakeInsuranceCreatePolicy.Name)
		assert.Equal(t, response.QuotationID, fakeInsuranceCreatePolicy.QuotationID)
		assert.Equal(t, response.DateOfBirth, fakeInsuranceCreatePolicy.DateOfBirth)
	})

	t.Run("Not should create a policy when have an invalid payload", func(t *testing.T) {
		defer clearAllDataBase()
		fakePartner := createAFakePartner()

		payload := map[string]interface{}{
			"sex": "A",
		}

		jsonData, err := json.Marshal(payload)
		assert.NoError(t, err)

		path := fmt.Sprintf("%s%s/policies", PartnerPath, fakePartner.ID)

		req, _ := http.NewRequest(http.MethodPost, path, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Not should create a policy when patner not exists", func(t *testing.T) {
		defer clearAllDataBase()
		createAFakePartner()

		payload := map[string]interface{}{
			"quotation_id":  fakeInsuranceCreatePolicy.QuotationID,
			"name":          fakeInsuranceCreatePolicy.Name,
			"sex":           fakeInsuranceCreatePolicy.Sex,
			"date_of_birth": fakeInsuranceCreatePolicy.DateOfBirth,
		}

		jsonData, err := json.Marshal(payload)
		assert.NoError(t, err)

		path := fmt.Sprintf("%s%s/policies", PartnerPath, "67e1a0349d00cb473900ac08")

		req, _ := http.NewRequest(http.MethodPost, path, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Not should create a policy when provider was error", func(t *testing.T) {
		defer clearAllDataBase()

		payload := map[string]interface{}{
			"quotation_id":  fakeInsuranceCreatePolicy.QuotationID,
			"name":          fakeInsuranceCreatePolicy.Name,
			"sex":           fakeInsuranceCreatePolicy.Sex,
			"date_of_birth": fakeInsuranceCreatePolicy.DateOfBirth,
		}

		fakePartner := createAFakePartner()
		setResponseInsurancePolicyError()

		jsonData, err := json.Marshal(payload)
		assert.NoError(t, err)

		path := fmt.Sprintf("%s%s/policies", PartnerPath, fakePartner.ID)

		req, _ := http.NewRequest(http.MethodPost, path, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

}

func TestGetPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, server, cleanUp, clearAllDataBase := testContext(ctrl)
	defer cleanUp()

	fakeInsuranceGetPolicy := partnerDomain.InsuranceProviderCreatePolicyResponse{
		ID:          uuid.New(),
		QuotationID: uuid.New(),
		Name:        "policy-test",
		Sex:         "F",
		DateOfBirth: "1998-09-28",
	}

	t.Run("Should return a sucessfuly a policy", func(t *testing.T) {
		defer clearAllDataBase()

		fakePartner := createAFakePartner()
		fakePolicy := createAFakePolicy(fakePartner.ID)
		setResponseGetInsurancePolicy(fakeInsuranceGetPolicy)

		path := fmt.Sprintf("%s%s/policies/%s", PartnerPath, fakePartner.ID, fakePolicy.ID)

		req, _ := http.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		var response partnersHandler.CreatePolicyResponseData
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, response.ID, fakePolicy.ID)
		assert.Equal(t, response.Sex, fakeInsuranceGetPolicy.Sex)
		assert.Equal(t, response.Name, fakeInsuranceGetPolicy.Name)
		assert.Equal(t, response.QuotationID, fakeInsuranceGetPolicy.QuotationID)
		assert.Equal(t, response.DateOfBirth, fakeInsuranceGetPolicy.DateOfBirth)
	})

	t.Run("Should return not found when pattern not bind to policy when user was try to get", func(t *testing.T) {
		defer clearAllDataBase()

		fakePartner := createAFakePartner()
		fakePolicy := createAFakePolicy(uuid.NewString())

		path := fmt.Sprintf("%s%s/policies/%s", PartnerPath, fakePartner.ID, fakePolicy.ID)

		req, _ := http.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		defer resp.Body.Close()
	})
}
