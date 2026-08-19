package validators

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validate is the shared validator instance reused across handlers.
var Validate = validator.New()

// ArticleRequest is the input body for create and update endpoints.
type ArticleRequest struct {
	Title    string `json:"title"    validate:"required,min=20"`
	Content  string `json:"content"  validate:"required,min=200"`
	Category string `json:"category" validate:"required,min=3"`
	Status   string `json:"status"   validate:"required,oneof=publish draft trash"`
}

// FormatErrors converts validator.ValidationErrors into a map of field → human-readable message.
func FormatErrors(err error) map[string]string {
	errs := make(map[string]string)
	for _, fe := range err.(validator.ValidationErrors) {
		field := fe.Field()
		switch fe.Tag() {
		case "required":
			errs[field] = field + " is required"
		case "min":
			errs[field] = fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
		case "oneof":
			errs[field] = fmt.Sprintf("%s must be one of: publish, draft, trash", field)
		default:
			errs[field] = fe.Error()
		}
	}
	return errs
}
