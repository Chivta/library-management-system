package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationErrors(err error) map[string]interface{} {
	var errors []ValidationError
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			errors = append(errors, ValidationError{
				Field:   fieldErr.Field(),
				Message: getErrorMessage(fieldErr),
			})
		}
	} else {
		errors = append(errors, ValidationError{
			Field:   "general",
			Message: err.Error(),
		})
	}

	return map[string]interface{}{
		"errors": errors,
	}
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}
