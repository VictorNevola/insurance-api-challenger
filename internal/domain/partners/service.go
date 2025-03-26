package partners

import (
	"context"
)

type (
	Service interface {
		CreatePartner(ctx context.Context, partner *PartnerEntity) (*PartnerEntity, error)
		CreateQuote(ctx context.Context, quote *QuoteEntity) (*QuoteEntity, error)
		CreatePolicy(ctx context.Context, policy *PolicyEntity) (*PolicyEntity, error)
		GetPolicy(ctx context.Context, partnerID, policyID string) (*PolicyEntity, error)
	}

	Servicer struct {
		partnerRepo       PartnerRepository
		quoteRepo         QuotesRepository
		policyRepo        PoliciesRepository
		insuranceProvider InsuranceProvider
	}

	ServiceParams struct {
		PartnerRepo             PartnerRepository
		QuoteRepo               QuotesRepository
		PolicyRepo              PoliciesRepository
		InsuranceClientProvider InsuranceProvider
	}
)

func NewService(data ServiceParams) *Servicer {
	return &Servicer{
		partnerRepo:       data.PartnerRepo,
		quoteRepo:         data.QuoteRepo,
		policyRepo:        data.PolicyRepo,
		insuranceProvider: data.InsuranceClientProvider,
	}
}

func (s *Servicer) CreatePartner(ctx context.Context, partner *PartnerEntity) (*PartnerEntity, error) {
	partnerExists, err := s.partnerRepo.GetByFilter(ctx, map[string]interface{}{"cnpj": partner.Cnpj})
	if err != nil {
		return nil, err
	}

	if partnerExists != nil {
		return nil, ErrPartnerAlreadyExists
	}

	err = s.partnerRepo.Create(ctx, partner)
	if err != nil {
		return nil, err
	}

	return partner, nil
}

func (s *Servicer) CreateQuote(ctx context.Context, quote *QuoteEntity) (*QuoteEntity, error) {
	partner, err := s.partnerRepo.GetByID(ctx, quote.PartnerID)
	if err != nil {
		return nil, err
	}

	if partner == nil {
		return nil, ErrPartnerNotFound
	}

	response, err := s.insuranceProvider.CreateQuotation(ctx, InsuranceProviderCreateQuotationRequest{
		Age: quote.Age,
		Sex: quote.Sex,
	})
	if err != nil {
		return nil, err
	}

	quoteCreated := &QuoteEntity{
		ProviderID: response.ProviderID,
		Age:        response.Age,
		Sex:        response.Sex,
		Price:      response.Price,
		PartnerID:  quote.PartnerID,
		CreatedAt:  quote.CreatedAt,
	}

	err = quoteCreated.ParseDateToEndOfDay(response.ExpiresAt)
	if err != nil {
		return nil, err
	}

	err = s.quoteRepo.Create(ctx, quoteCreated)
	if err != nil {
		return nil, err
	}

	return quoteCreated, nil
}

func (s *Servicer) CreatePolicy(ctx context.Context, policy *PolicyEntity) (*PolicyEntity, error) {
	partner, err := s.partnerRepo.GetByID(ctx, policy.PartnerID)
	if err != nil {
		return nil, err
	}

	if partner == nil {
		return nil, ErrPartnerNotFound
	}

	response, err := s.insuranceProvider.CreatePolicy(ctx, InsuranceProviderCreatePolicyRequest{
		QuotationID: policy.QuotationID,
		Name:        policy.Name,
		Sex:         string(policy.Sex),
		DateOfBirth: policy.DateOfBirth,
	})
	if err != nil {
		return nil, err
	}

	policy.ProviderID = response.ID
	err = s.policyRepo.Create(ctx, policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (s *Servicer) GetPolicy(ctx context.Context, partnerID, policyID string) (*PolicyEntity, error) {
	partner, err := s.partnerRepo.GetByID(ctx, partnerID)
	if err != nil {
		return nil, err
	}

	if partner == nil {
		return nil, ErrPartnerNotFound
	}

	userHasPolicy, err := s.policyRepo.GetByIdAndPartnerID(ctx, policyID, partnerID)
	if err != nil {
		return nil, err
	}

	if userHasPolicy == nil {
		return nil, ErrPolicyNotFound
	}

	policy, err := s.insuranceProvider.GetPolicy(ctx, userHasPolicy.ProviderID.String())
	if err != nil {
		return nil, err
	}

	return &PolicyEntity{
		ID:          userHasPolicy.ID,
		Sex:         SexEnum(policy.Sex),
		Name:        policy.Name,
		QuotationID: policy.QuotationID,
		DateOfBirth: policy.DateOfBirth,
	}, nil
}
