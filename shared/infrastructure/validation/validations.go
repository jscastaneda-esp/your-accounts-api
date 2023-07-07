package validation

import (
	"log"
	"your-accounts-api/shared/domain/validation"

	"github.com/gofiber/fiber/v2"
)

func Validate(c *fiber.Ctx, request any) bool {
	if err := c.BodyParser(request); err != nil {
		log.Println("Error request body parser:", err)
		c.Status(fiber.StatusBadRequest)
		return false
	}

	if errors := validation.ValidateStruct(request); errors != nil {
		c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
		return false
	}

	return true
}
