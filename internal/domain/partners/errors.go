package partners

import (
	"github.com/gofiber/fiber/v2"
)

var (
	ErrPartnerAlreadyExists = fiber.NewError(fiber.StatusConflict, "partner already exists")
	ErrPartnerNotFound      = fiber.NewError(fiber.StatusNotFound, "partner not found")
	ErrPolicyNotFound       = fiber.NewError(fiber.StatusNotFound, "policy not found")
)
