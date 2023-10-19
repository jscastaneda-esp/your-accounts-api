package validation

import (
	"your-accounts-api/shared/domain/validation"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

func Validate(c *fiber.Ctx, request any) bool {
	if err := c.BodyParser(request); err != nil {
		log.Error("Error request body parser:", err)
		c.Status(fiber.StatusBadRequest)
		return false
	}

	if errors := validation.ValidateStruct(request); errors != nil {
		c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
		return false
	}

	return true
}

func ValidateSlice(c *fiber.Ctx, request any, constraint string) bool {
	if err := c.BodyParser(request); err != nil {
		log.Error("Error request body parser:", err)
		c.Status(fiber.StatusBadRequest)
		return false
	}

	if errors := validation.ValidateVariable(request, constraint); errors != nil {
		c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
		return false
	}

	return true
}
