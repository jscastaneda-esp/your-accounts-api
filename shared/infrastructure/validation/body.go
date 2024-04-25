package validation

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const RequestBody = "body"

func RequestBodyValid(t any) fiber.Handler {
	typ := reflect.TypeOf(t)

	return func(c *fiber.Ctx) error {
		instance := reflect.New(typ).Interface()

		if err := c.BodyParser(instance); err != nil {
			log.Error("Error request body parser:", err)
			return fiber.ErrBadRequest
		}

		if err := validate.Struct(instance); err != nil {
			log.Error("Error request body validation:", err)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(getErrors(err))
		}

		c.Locals(RequestBody, instance)
		return c.Next()
	}
}
