// TODO: Pendientes tests

package controller

import (
	"errors"
	"log"
	"strings"

	"api-your-accounts/shared/infrastructure/db"
	"api-your-accounts/shared/infrastructure/validation"
	"api-your-accounts/user/application"
	"api-your-accounts/user/domain"
	"api-your-accounts/user/infrastructure/model"
	"api-your-accounts/user/infrastructure/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateUserHandler(c *fiber.Ctx) error {
	request := new(model.CreateRequest)
	if err := c.BodyParser(request); err != nil {
		log.Println("Error request body parser:", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if errors := validation.ValidateStruct(request); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	repo := repository.NewGORMRepository(db.DB)
	user := &domain.User{
		UUID:  request.UUID,
		Email: request.Email,
	}

	exists, err := application.Exists(repo, c.UserContext(), user.UUID, user.Email)
	if exists {
		return fiber.NewError(fiber.StatusConflict, "User already exists")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Error sign up user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error sign up user")
	}

	result, err := application.SignUp(repo, c.UserContext(), user)
	if err != nil {
		log.Println("Error sign up user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error sign up user")
	}

	return c.JSON(model.CreateResponse{
		ID:        result.ID,
		UUID:      result.UUID,
		Email:     result.Email,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	})
}

func LoginHandler(c *fiber.Ctx) error {
	request := new(model.LoginRequest)
	if err := c.BodyParser(request); err != nil {
		log.Println("Error request body parser:", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if errors := validation.ValidateStruct(request); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	repo := repository.NewGORMRepository(db.DB)
	token, err := application.Login(repo, c.UserContext(), request.UUID, strings.ToLower(request.Email))
	if err != nil {
		log.Println("Error login user:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Error login user")
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
