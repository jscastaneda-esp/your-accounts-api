package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Field      string `json:"field"`
	Constraint string `json:"constraint"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func getErrors(err error) []ErrorResponse {
	var errors []ErrorResponse
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			error := ErrorResponse{
				Field: err.Namespace(),
			}

			if err.Param() == "" {
				error.Constraint = err.Tag()
			} else {
				error.Constraint = fmt.Sprintf("%s=%s", err.Tag(), err.Param())
			}

			errors = append(errors, error)
		}
	}

	return errors
}
