package validator

import (
	"github.com/go-playground/validator/v10"
)

var NewVal = validator.New()

// RegisterValidators registers the validators
func RegisterCustomValidators() {
	NewVal.RegisterValidation("project_status_enum", projectStatusEnum)
}

func projectStatusEnum(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "backlog", "developing", "done":
		return true
	default:
		return false
	}
}
