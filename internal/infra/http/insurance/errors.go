package insurance

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type (
	apiError struct {
		Message string `json:"message"`
	}
)

func handlerErrors(statusCode int, responseBody []byte) error {
	var apiError apiError

	if statusCode >= 400 && statusCode < 500 {
		err := json.Unmarshal(responseBody, &apiError)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("unexpected client error: %s", string(responseBody)))
		}

		return fiber.NewError(fiber.StatusBadRequest, apiError.Message)
	}

	return nil
}
