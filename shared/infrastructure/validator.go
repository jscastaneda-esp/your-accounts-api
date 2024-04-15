package infrastructure

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Field      string `json:"field"`
	Constraint string `json:"constraint"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (v *CustomValidator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		errors := getErrors(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, errors)
	}

	return nil
}

func (v *CustomValidator) ValidateVariable(i interface{}, tag string) error {
	if err := v.validator.Var(i, tag); err != nil {
		errors := getErrors(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, errors)
	}

	return nil
}

func getErrors(err error) []*ErrorResponse {
	var errors []*ErrorResponse
	for _, err := range err.(validator.ValidationErrors) {
		error := &ErrorResponse{
			Field: err.Field(),
		}

		if err.Param() == "" {
			error.Constraint = err.Tag()
		} else {
			error.Constraint = fmt.Sprintf("%s=%s", err.Tag(), err.Param())
		}

		errors = append(errors, error)
	}

	return errors
}

func NewCustomValidator() *CustomValidator {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{
		validator: validate,
	}
}
