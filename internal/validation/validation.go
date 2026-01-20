package validation

import (
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	v := validator.New(validator.WithRequiredStructEnabled())

	Validate = v
}
