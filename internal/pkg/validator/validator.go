package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

var validate = validator.New()

func BodyData(data interface{}) *fiber.Error {
	validationErrors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return FormatErrors(validationErrors)
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if fiberErr, ok := err.(*fiber.Error); ok {
		return c.Status(fiberErr.Code).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: fiberErr.Message,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(GlobalErrorHandlerResp{
		Success: false,
		Message: "internal server error",
	})
}

func FormatErrors(errs []ErrorResponse) *fiber.Error {
	errMsgs := make([]string, len(errs))

	for index, err := range errs {
		errMsgs[index] = fmt.Sprintf(
			"[%s]: '%v' | Needs to implement '%s'",
			err.FailedField,
			err.Value,
			err.Tag,
		)
	}

	if len(errMsgs) > 0 {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, ", "),
		}
	}

	return nil
}
