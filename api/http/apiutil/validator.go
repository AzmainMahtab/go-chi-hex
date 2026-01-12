package apiutil

import (
	"github.com/AzmainMahtab/go-chi-hex/pkg/jsonutil"
	"github.com/go-playground/validator/v10"
)

// Global validator instance
var v = validator.New()

// ValidateStruct is a helper to run validation and return domain-friendly errors
func ValidateStruct(s any) []jsonutil.ErrorItem {
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	var errorItems []jsonutil.ErrorItem
	// Cast to validator.ValidationErrors to loop through fields
	for _, err := range err.(validator.ValidationErrors) {
		errorItems = append(errorItems, jsonutil.ErrorItem{
			Field:   err.Field(),
			Message: getCustomMessage(err),
		})
	}
	return errorItems
}

func getCustomMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "e164":
		return "Invalid international phone format"
	default:
		return "Invalid value"
	}
}
