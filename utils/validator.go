package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// NewValidator func for create a new validator for model fields.
//
//	validate := utils.NewValidator()
//	if err := validate.Struct(object); err != nil { }
func NewValidator() *validator.Validate {
	validate := validator.New()

	// Custom validation for uuid.UUID fields.
	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		_, err := uuid.Parse(field)
		return err == nil
	})

	return validate
}

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) map[string]string {
	// Define fields map.
	fields := map[string]string{}
	if err == nil {
		return fields
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fields
	}

	// Make error message for each invalid field.
	for _, err := range validationErrors {
		fields[err.Field()] = err.Error()
	}

	return fields
}
