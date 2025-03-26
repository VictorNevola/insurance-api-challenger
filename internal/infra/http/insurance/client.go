package insurance

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main-api/internal/domain/partners"
	"main-api/internal/infra/cache"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type (
	InsuranceProviderClient struct {
		baseURL      string
		apiKey       string
		cacheStorage cache.CacheStore
		Client       *http.Client
	}

	authenticateResponse struct {
		AcessToken string `json:"access_token"`
	}

	createQuotationResponse struct {
		ID        uuid.UUID `json:"id"`
		Sex       string    `json:"sex"`
		ExpiresAt string    `json:"expire_at"`
		Age       uint      `json:"age"`
		Price     float64   `json:"price"`
	}

	policyResponse struct {
		ID          uuid.UUID `json:"id"`
		QuotationID uuid.UUID `json:"quotation_id"`
		DateBirth   string    `json:"date_of_birth"`
		Name        string    `json:"name"`
		Sex         string    `json:"sex"`
	}
)

var (
	jwtKey = "insurance-provider-jwt-token"
)

func NewInsuranceProviderClient(cache cache.CacheStore, baseURL, apiKey string) *InsuranceProviderClient {

	return &InsuranceProviderClient{
		baseURL:      baseURL,
		cacheStorage: cache,
		apiKey:       apiKey,
		Client:       &http.Client{},
	}
}

func (i *InsuranceProviderClient) Authenticate() (*authenticateResponse, error) {
	req, err := http.NewRequest("POST", i.baseURL+"/auth", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", i.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var authResponse authenticateResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return nil, err
	}

	return &authResponse, nil
}

func (i *InsuranceProviderClient) CreateQuotation(
	ctx context.Context,
	data partners.InsuranceProviderCreateQuotationRequest,
) (*partners.InsuranceProviderCreateQuotationResponse, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	response, err := i.doRequestWithAuth(ctx, "POST", "quotations", jsonData)
	if err != nil {
		return nil, err
	}

	var createQuotationResponse createQuotationResponse
	err = json.Unmarshal(response, &createQuotationResponse)
	if err != nil {
		return nil, err
	}

	return &partners.InsuranceProviderCreateQuotationResponse{
		ProviderID: createQuotationResponse.ID,
		Age:        createQuotationResponse.Age,
		Price:      createQuotationResponse.Price,
		ExpiresAt:  createQuotationResponse.ExpiresAt,
		Sex:        partners.SexEnum(createQuotationResponse.Sex),
	}, nil
}

func (i *InsuranceProviderClient) CreatePolicy(
	ctx context.Context,
	data partners.InsuranceProviderCreatePolicyRequest,
) (*partners.InsuranceProviderCreatePolicyResponse, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	body, err := i.doRequestWithAuth(ctx, "POST", "policies", jsonData)
	if err != nil {
		return nil, err
	}

	var createPolicyResponse policyResponse
	err = json.Unmarshal(body, &createPolicyResponse)
	if err != nil {
		return nil, err
	}

	return &partners.InsuranceProviderCreatePolicyResponse{
		ID:          createPolicyResponse.ID,
		QuotationID: createPolicyResponse.QuotationID,
		Name:        createPolicyResponse.Name,
		Sex:         createPolicyResponse.Sex,
		DateOfBirth: createPolicyResponse.DateBirth,
	}, nil
}

func (i *InsuranceProviderClient) GetPolicy(ctx context.Context, policyID string) (*partners.InsuranceProviderCreatePolicyResponse, error) {
	body, err := i.doRequestWithAuth(ctx, "GET", "policies/"+policyID, nil)
	if err != nil {
		return nil, err
	}

	var response policyResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &partners.InsuranceProviderCreatePolicyResponse{
		ID:          response.ID,
		QuotationID: response.QuotationID,
		Name:        response.Name,
		DateOfBirth: response.DateBirth,
		Sex:         response.Sex,
	}, nil
}

func (i *InsuranceProviderClient) getToken(ctx context.Context) (string, error) {
	tokenInCache, err := i.cacheStorage.Get(ctx, jwtKey)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		return "", err
	}

	if tokenInCache != "" {
		return tokenInCache, nil
	}

	authResponse, err := i.Authenticate()
	if err != nil {
		return "", err
	}

	err = i.cacheStorage.Set(ctx, jwtKey, authResponse.AcessToken, 10*time.Minute)
	if err != nil {
		return "", err
	}

	return authResponse.AcessToken, nil
}

func (i *InsuranceProviderClient) doRequestWithAuth(ctx context.Context, method, url string, payload []byte) ([]byte, error) {
	token, err := i.getToken(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s/%s", i.baseURL, url),
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if userErr := handlerErrors(resp.StatusCode, body); userErr != nil {
		return nil, userErr
	}

	return body, nil
}
