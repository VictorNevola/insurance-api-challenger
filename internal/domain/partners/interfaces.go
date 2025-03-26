package partners

import (
	"context"

	"github.com/google/uuid"
)

type (
	PartnerRepository interface {
		GetByFilter(ctx context.Context, filter map[string]interface{}) (*PartnerEntity, error)
		GetByID(ctx context.Context, id string) (*PartnerEntity, error)
		Create(ctx context.Context, partner *PartnerEntity) error
	}

	QuotesRepository interface {
		Create(ctx context.Context, quote *QuoteEntity) error
	}

	PoliciesRepository interface {
		Create(ctx context.Context, policy *PolicyEntity) error
		GetByIdAndPartnerID(ctx context.Context, policyID, partnerID string) (*PolicyEntity, error)
	}

	InsuranceProviderCreateQuotationRequest struct {
		Age uint
		Sex SexEnum
	}

	InsuranceProviderCreateQuotationResponse struct {
		ProviderID uuid.UUID
		Age        uint
		Price      float64
		Sex        SexEnum
		ExpiresAt  string
	}

	InsuranceProviderCreatePolicyRequest struct {
		QuotationID uuid.UUID `json:"quotation_id"`
		Name        string    `json:"name"`
		Sex         string    `json:"sex"`
		DateOfBirth string    `json:"date_of_birth"`
	}

	InsuranceProviderCreatePolicyResponse struct {
		ID          uuid.UUID
		QuotationID uuid.UUID
		Name        string
		Sex         string
		DateOfBirth string
	}

	InsuranceProvider interface {
		CreateQuotation(
			ctx context.Context,
			data InsuranceProviderCreateQuotationRequest,
		) (*InsuranceProviderCreateQuotationResponse, error)
		CreatePolicy(
			ctx context.Context,
			data InsuranceProviderCreatePolicyRequest,
		) (*InsuranceProviderCreatePolicyResponse, error)
		GetPolicy(ctx context.Context, policyID string) (*InsuranceProviderCreatePolicyResponse, error)
	}
)
