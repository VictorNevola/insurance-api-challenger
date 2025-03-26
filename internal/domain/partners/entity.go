package partners

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type (
	SexEnum string

	PartnerEntity struct {
		ID        string
		Name      string
		Cnpj      string
		CreatedAt time.Time
	}

	QuoteEntity struct {
		ID         string
		ProviderID uuid.UUID
		Age        uint
		Sex        SexEnum
		PartnerID  string
		Price      float64
		ExpiresAt  time.Time
		CreatedAt  time.Time
	}

	PolicyEntity struct {
		ID          string
		QuotationID uuid.UUID
		ProviderID  uuid.UUID
		Sex         SexEnum
		Name        string
		DateOfBirth string
		PartnerID   string
	}
)

const (
	SexMale    SexEnum = "M"
	SexFemale  SexEnum = "F"
	SexNeutral SexEnum = "N"
)

func NewEntity(name, cnpj string) *PartnerEntity {
	return &PartnerEntity{
		Name:      name,
		Cnpj:      cnpj,
		CreatedAt: time.Now(),
	}
}

func (e *PartnerEntity) SetID(id string) error {
	e.ID = id

	return nil
}

func NewQuoteEntity(age uint, sex string, partnerID string) *QuoteEntity {
	return &QuoteEntity{
		Age:       age,
		Sex:       SexEnum(strings.ToUpper(sex)),
		PartnerID: partnerID,
		CreatedAt: time.Now(),
	}
}

func (e *QuoteEntity) ParseDateToEndOfDay(dateStr string) error {
	const layout = "2006-01-02"

	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return err
	}

	endOfDay := time.Date(
		parsedDate.Year(),
		parsedDate.Month(),
		parsedDate.Day(),
		23, 59, 59, 0,
		parsedDate.Location(),
	)

	e.ExpiresAt = endOfDay

	return nil
}
