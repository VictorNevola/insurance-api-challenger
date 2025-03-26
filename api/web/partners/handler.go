package partners

import (
	"main-api/internal/domain/partners"
	"main-api/internal/pkg/validator"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type (
	HTTPHandler struct {
		service partners.Service
	}

	CreatePartnerRequestData struct {
		Name string `json:"name" validate:"required,min=3,max=255"`
		Cnpj string `json:"cnpj" validate:"required,min=14,max=14"`
	}

	CreatePartnerResponseData struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Cnpj      string    `json:"cnpj"`
		CreatedAt time.Time `json:"created_at"`
	}

	CreateQuoteData struct {
		Age uint   `json:"age" validate:"required,min=0,max=99"`
		Sex string `json:"sex" validate:"required,oneof=m M f F n N"`
	}

	CreateQuoteResponseData struct {
		ID        string    `json:"id"`
		Age       uint      `json:"age"`
		Sex       string    `json:"sex"`
		Price     float64   `json:"price"`
		ExpiresAt time.Time `json:"expires_at"`
		CreatedAt time.Time `json:"created_at"`
	}

	CreatePolicyData struct {
		QuotationID uuid.UUID `json:"quotation_id" validate:"required"`
		Name        string    `json:"name" validate:"required,min=3,max=255"`
		Sex         string    `json:"sex" validate:"required,oneof=m M f F n N"`
		DateOfBirth string    `json:"date_of_birth" validate:"required"`
	}

	CreatePolicyResponseData struct {
		ID          string    `json:"id"`
		Sex         string    `json:"sex"`
		Name        string    `json:"name"`
		QuotationID uuid.UUID `json:"quotation_id"`
		DateOfBirth string    `json:"date_of_birth"`
	}
)

func NewHTTPHandler(app *fiber.App, service partners.Service) {
	httpHandler := HTTPHandler{
		service: service,
	}

	app.Route("/partners", func(r fiber.Router) {
		r.Post("/", httpHandler.CreatePartner)
		r.Post("/:partner_id/quotes", httpHandler.CreateQuote)
		r.Post("/:partner_id/policies", httpHandler.CreatePolicy)
		r.Get("/:partner_id/policies/:policy_id", httpHandler.GetPolicy)
	})
}

func (h *HTTPHandler) CreatePartner(c *fiber.Ctx) error {
	bodyData := new(CreatePartnerRequestData)
	if err := c.BodyParser(bodyData); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	if err := validator.BodyData(bodyData); err != nil {
		return err
	}

	partnerEntity := partners.NewEntity(bodyData.Name, bodyData.Cnpj)
	partner, err := h.service.CreatePartner(c.Context(), partnerEntity)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(CreatePartnerResponseData{
		ID:        partner.ID,
		Name:      partner.Name,
		Cnpj:      partner.Cnpj,
		CreatedAt: partner.CreatedAt,
	})
}

func (h *HTTPHandler) CreateQuote(c *fiber.Ctx) error {
	bodyData := new(CreateQuoteData)
	if err := c.BodyParser(bodyData); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	if err := validator.BodyData(bodyData); err != nil {
		return err
	}

	result, err := h.service.CreateQuote(c.Context(), partners.NewQuoteEntity(
		bodyData.Age,
		bodyData.Sex,
		c.Params("partner_id"),
	))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(CreateQuoteResponseData{
		ID:        result.ProviderID.String(),
		Age:       result.Age,
		Sex:       string(result.Sex),
		Price:     result.Price,
		ExpiresAt: result.ExpiresAt,
		CreatedAt: result.CreatedAt,
	})
}

func (h HTTPHandler) CreatePolicy(c *fiber.Ctx) error {
	bodyData := new(CreatePolicyData)
	if err := c.BodyParser(bodyData); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	if err := validator.BodyData(bodyData); err != nil {
		return err
	}

	result, err := h.service.CreatePolicy(c.Context(), &partners.PolicyEntity{
		QuotationID: bodyData.QuotationID,
		Sex:         partners.SexEnum(bodyData.Sex),
		Name:        bodyData.Name,
		DateOfBirth: bodyData.DateOfBirth,
		PartnerID:   c.Params("partner_id"),
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(CreatePolicyResponseData{
		ID:          result.ID,
		Sex:         string(result.Sex),
		Name:        result.Name,
		QuotationID: result.QuotationID,
		DateOfBirth: result.DateOfBirth,
	})
}

func (h HTTPHandler) GetPolicy(c *fiber.Ctx) error {
	policy, err := h.service.GetPolicy(c.Context(), c.Params("partner_id"), c.Params("policy_id"))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(CreatePolicyResponseData{
		ID:          policy.ID,
		Sex:         string(policy.Sex),
		Name:        policy.Name,
		QuotationID: policy.QuotationID,
		DateOfBirth: policy.DateOfBirth,
	})
}
