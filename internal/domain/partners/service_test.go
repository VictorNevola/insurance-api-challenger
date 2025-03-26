package partners_test

import (
	"context"
	"errors"
	"main-api/internal/domain/partners"
	mocks "main-api/internal/infra/repository/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestServiceCreatePartner(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	partnersRepo := mocks.NewMockPartnerRepository(ctrl)

	service := partners.NewService(partners.ServiceParams{
		PartnerRepo: partnersRepo,
	})

	t.Run("Should return success when creating a partner", func(t *testing.T) {
		partner := partners.NewEntity("180 Seguros", "12345678901234")

		partnersRepo.EXPECT().GetByFilter(gomock.Any(), gomock.Any()).Return(nil, nil)
		partnersRepo.EXPECT().Create(gomock.Any(), partner).Return(nil)

		partnerCreated, err := service.CreatePartner(t.Context(), partner)

		assert.Nil(t, err)
		assert.Equal(t, partner, partnerCreated)
	})

	t.Run("Should return error when partner already exists", func(t *testing.T) {
		partner := partners.NewEntity("180 Seguros", "12345678901234")

		partnersRepo.EXPECT().GetByFilter(gomock.Any(), gomock.Any()).Return(partner, nil)

		partnerCreated, err := service.CreatePartner(t.Context(), partner)

		assert.Nil(t, partnerCreated)
		assert.Equal(t, partners.ErrPartnerAlreadyExists, err)
	})

	t.Run("Should return error when creating a partner and have an internal error", func(t *testing.T) {
		partner := partners.NewEntity("180 Seguros", "12345678901234")

		partnersRepo.EXPECT().GetByFilter(gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error"))

		partnerCreated, err := service.CreatePartner(t.Context(), partner)

		assert.Nil(t, partnerCreated)
		assert.Equal(t, "internal error", err.Error())
	})
}

func TestServiceCreateQuote(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	partnersRepo := mocks.NewMockPartnerRepository(ctrl)
	quotesRepo := mocks.NewMockQuotesRepository(ctrl)
	insuranceProviderClient := mocks.NewMockInsuranceProvider(ctrl)

	service := partners.NewService(partners.ServiceParams{
		PartnerRepo:             partnersRepo,
		QuoteRepo:               quotesRepo,
		InsuranceClientProvider: insuranceProviderClient,
	})

	fakePartner := partners.PartnerEntity{
		ID:        uuid.NewString(),
		Name:      "partner-test",
		Cnpj:      "12345678901234",
		CreatedAt: time.Now(),
	}

	insuranceProviderFakeRes := partners.InsuranceProviderCreateQuotationResponse{
		ProviderID: uuid.New(),
		Age:        26,
		Price:      12.78,
		Sex:        "M",
		ExpiresAt:  "2999-12-31",
	}

	quote := &partners.QuoteEntity{
		Age:       26,
		Sex:       "M",
		PartnerID: fakePartner.ID,
	}

	t.Run("Should return success when creating a quote", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
			Return(&fakePartner, nil)
		insuranceProviderClient.EXPECT().CreateQuotation(gomock.Any(), gomock.Any()).
			Return(&insuranceProviderFakeRes, nil)
		quotesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, quote *partners.QuoteEntity) error {
				quote.ID = uuid.NewString()
				return nil
			},
		)

		createdQuote, err := service.CreateQuote(t.Context(), quote)

		assert.Nil(t, err)
		assert.NotNil(t, createdQuote)
		assert.Equal(t, createdQuote.ProviderID, insuranceProviderFakeRes.ProviderID)
		assert.Equal(t, createdQuote.Price, insuranceProviderFakeRes.Price)
		assert.Equal(t, createdQuote.Age, quote.Age)
		assert.Equal(t, createdQuote.Sex, quote.Sex)
		assert.Equal(t, createdQuote.PartnerID, fakePartner.ID)
	})

	t.Run("Should return error when partner is not found", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
			Return(nil, nil)

		createdQuote, err := service.CreateQuote(context.Background(), quote)

		assert.Nil(t, createdQuote)
		assert.Equal(t, partners.ErrPartnerNotFound, err)
	})

	t.Run("Should return error when GetByID fails", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), quote.PartnerID).Return(nil, errors.New("database error"))

		createdQuote, err := service.CreateQuote(context.Background(), quote)

		assert.Nil(t, createdQuote)
		assert.EqualError(t, err, "database error")
	})

	t.Run("Should return error when CreateQuotation fails", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), quote.PartnerID).Return(&fakePartner, nil)
		insuranceProviderClient.EXPECT().CreateQuotation(gomock.Any(), gomock.Any()).Return(nil, errors.New("quotation error"))

		createdQuote, err := service.CreateQuote(context.Background(), quote)

		assert.Nil(t, createdQuote)
		assert.EqualError(t, err, "quotation error")
	})

	t.Run("Should return error when ParseDateToEndOfDay fails", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), quote.PartnerID).Return(&fakePartner, nil)
		insuranceProviderClient.EXPECT().CreateQuotation(gomock.Any(), gomock.Any()).Return(&partners.InsuranceProviderCreateQuotationResponse{
			ProviderID: uuid.New(),
			Age:        26,
			Price:      12.78,
			Sex:        "M",
			ExpiresAt:  "invalid-date", // Data inválida para forçar o erro
		}, nil)

		createdQuote, err := service.CreateQuote(context.Background(), quote)

		assert.Nil(t, createdQuote)
		assert.Contains(t, err.Error(), "parsing time") // Verifica se o erro é de parsing
	})

	t.Run("Should return error when QuoteRepo.Create fails", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), quote.PartnerID).Return(&fakePartner, nil)
		insuranceProviderClient.EXPECT().CreateQuotation(gomock.Any(), gomock.Any()).Return(&partners.InsuranceProviderCreateQuotationResponse{
			ProviderID: uuid.New(),
			Age:        26,
			Price:      12.78,
			Sex:        "M",
			ExpiresAt:  time.Now().Format("2006-01-02"),
		}, nil)
		quotesRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("repository error"))

		createdQuote, err := service.CreateQuote(context.Background(), quote)

		assert.Nil(t, createdQuote)
		assert.EqualError(t, err, "repository error")
	})
}

func TestServiceCreatePolicy(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	partnersRepo := mocks.NewMockPartnerRepository(ctrl)
	quotesRepo := mocks.NewMockQuotesRepository(ctrl)
	policyRepo := mocks.NewMockPoliciesRepository(ctrl)
	insuranceProviderClient := mocks.NewMockInsuranceProvider(ctrl)

	service := partners.NewService(partners.ServiceParams{
		PartnerRepo:             partnersRepo,
		QuoteRepo:               quotesRepo,
		PolicyRepo:              policyRepo,
		InsuranceClientProvider: insuranceProviderClient,
	})

	fakePartner := partners.PartnerEntity{
		ID:        uuid.NewString(),
		Name:      "partner-test",
		Cnpj:      "12345678901234",
		CreatedAt: time.Now(),
	}

	insuranceProviderFakeRes := partners.InsuranceProviderCreatePolicyResponse{
		ID:          uuid.New(),
		QuotationID: uuid.New(),
		Name:        "policy-test",
		Sex:         "F",
		DateOfBirth: "1998-09-28",
	}

	t.Run("Should create a policy succesfuly and return", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
			Return(&fakePartner, nil)
		insuranceProviderClient.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).
			Return(&insuranceProviderFakeRes, nil)
		policyRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, policy *partners.PolicyEntity) error {
				policy.ID = uuid.NewString()
				return nil
			},
		)

		policyCreated, err := service.CreatePolicy(t.Context(), &partners.PolicyEntity{
			QuotationID: uuid.New(),
			PartnerID:   uuid.New().String(),
			Sex:         "F",
			Name:        "policy-test",
			DateOfBirth: "1998-09-28",
		})

		assert.NoError(t, err)
		assert.NotNil(t, policyCreated)
		assert.Equal(t, policyCreated.ProviderID, insuranceProviderFakeRes.ID)
		assert.Equal(t, string(policyCreated.Sex), string(insuranceProviderFakeRes.Sex))
		assert.Equal(t, policyCreated.Name, insuranceProviderFakeRes.Name)
		assert.Equal(t, policyCreated.DateOfBirth, insuranceProviderFakeRes.DateOfBirth)
	})

	t.Run("Not should create a policy when patner not found", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
			Return(nil, nil)

		createdPolicy, err := service.CreatePolicy(context.Background(), &partners.PolicyEntity{
			QuotationID: uuid.New(),
			PartnerID:   uuid.New().String(),
			Sex:         "F",
			Name:        "policy-test",
			DateOfBirth: "1998-09-28",
		})

		assert.Nil(t, createdPolicy)
		assert.Equal(t, partners.ErrPartnerNotFound, err)
	})

	t.Run("Not should create a policy when provider was error", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
			Return(&fakePartner, nil)
		insuranceProviderClient.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).
			Return(nil, errors.New(`{"message": "The quotation was expired"}`))

		createdPolicy, err := service.CreatePolicy(context.Background(), &partners.PolicyEntity{
			QuotationID: uuid.New(),
			PartnerID:   uuid.New().String(),
			Sex:         "F",
			Name:        "policy-test",
			DateOfBirth: "1998-09-28",
		})

		assert.Nil(t, createdPolicy)
		assert.Equal(t, "{\"message\": \"The quotation was expired\"}", err.Error())
	})
}

func TestServiceGetPolicy(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	partnersRepo := mocks.NewMockPartnerRepository(ctrl)
	policyRepo := mocks.NewMockPoliciesRepository(ctrl)
	insuranceProviderClient := mocks.NewMockInsuranceProvider(ctrl)

	service := partners.NewService(partners.ServiceParams{
		PartnerRepo:             partnersRepo,
		PolicyRepo:              policyRepo,
		InsuranceClientProvider: insuranceProviderClient,
	})

	fakePartner := partners.PartnerEntity{
		ID:        uuid.NewString(),
		Name:      "partner-test",
		Cnpj:      "12345678901234",
		CreatedAt: time.Now(),
	}

	insuranceProviderFakeRes := partners.InsuranceProviderCreatePolicyResponse{
		ID:          uuid.New(),
		QuotationID: uuid.New(),
		Name:        "policy-test",
		Sex:         "F",
		DateOfBirth: "1998-09-28",
	}

	fakePolicyCreated := partners.PolicyEntity{
		ID:          uuid.NewString(),
		QuotationID: insuranceProviderFakeRes.QuotationID,
		PartnerID:   fakePartner.ID,
		ProviderID:  insuranceProviderFakeRes.ID,
		Sex:         partners.SexEnum(insuranceProviderFakeRes.Sex),
		Name:        insuranceProviderFakeRes.Name,
		DateOfBirth: insuranceProviderFakeRes.DateOfBirth,
	}

	t.Run("Should get a policy succesfuly and return", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
			Return(&fakePartner, nil)
		policyRepo.EXPECT().GetByIdAndPartnerID(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&fakePolicyCreated, nil)
		insuranceProviderClient.EXPECT().GetPolicy(gomock.Any(), gomock.Any()).
			Return(&insuranceProviderFakeRes, nil)

		policyCreated, err := service.GetPolicy(t.Context(), fakePartner.ID, fakePolicyCreated.ID)

		assert.NoError(t, err)
		assert.NotNil(t, policyCreated)
		assert.Equal(t, policyCreated.ID, fakePolicyCreated.ID)
		assert.Equal(t, string(policyCreated.Sex), string(insuranceProviderFakeRes.Sex))
		assert.Equal(t, policyCreated.Name, insuranceProviderFakeRes.Name)
		assert.Equal(t, policyCreated.DateOfBirth, insuranceProviderFakeRes.DateOfBirth)
	})

	t.Run("Not should return a policy when patner not found", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
			Return(nil, nil)

		policyCreated, err := service.GetPolicy(t.Context(), fakePartner.ID, fakePolicyCreated.ID)

		assert.Nil(t, policyCreated)
		assert.Equal(t, partners.ErrPartnerNotFound, err)
	})

	t.Run("Not should return a policy when partner not bind to this policy", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
			Return(&fakePartner, nil)
		policyRepo.EXPECT().GetByIdAndPartnerID(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil)

		policyCreated, err := service.GetPolicy(t.Context(), uuid.New().String(), fakePolicyCreated.ID)

		assert.Nil(t, policyCreated)
		assert.Equal(t, partners.ErrPolicyNotFound, err)
	})

	t.Run("Should return erro when try to get a policy and provider return error", func(t *testing.T) {
		partnersRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
			Return(&fakePartner, nil)
		policyRepo.EXPECT().GetByIdAndPartnerID(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&fakePolicyCreated, nil)
		insuranceProviderClient.EXPECT().GetPolicy(gomock.Any(), gomock.Any()).
			Return(nil, errors.New(`{"message": "The policy not found"}`))

		policyCreated, err := service.GetPolicy(t.Context(), fakePartner.ID, fakePolicyCreated.ID)
		assert.Nil(t, policyCreated)
		assert.Equal(t, "{\"message\": \"The policy not found\"}", err.Error())
	})
}
