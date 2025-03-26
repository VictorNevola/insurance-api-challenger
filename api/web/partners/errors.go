package partners

import (
	"main-api/internal/domain/partners"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrosMapped = map[error]*fiber.Error{
		partners.ErrPartnerAlreadyExists: fiber.NewError(
			fiber.StatusConflict,
			partners.ErrPartnerAlreadyExists.Error(),
		),
		partners.ErrPartnerNotFound: fiber.NewError(
			fiber.StatusNotFound,
			partners.ErrPartnerNotFound.Error(),
		),
	}
)

func HandlerCorrectlyErrorsStatus(status error) *fiber.Error {
	if err, ok := ErrosMapped[status]; ok {
		return err
	}

	return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
}
