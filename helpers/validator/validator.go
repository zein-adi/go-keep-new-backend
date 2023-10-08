package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"strings"
)

// New PlaygroundValidator Wrapper
func New() *Validator {
	return &Validator{
		playgroundValidator: validator.New(),
	}
}

type Validator struct {
	playgroundValidator *validator.Validate
}

func (r *Validator) ValidateMap(data map[string]interface{}, rules map[string]interface{}) error {
	r.playgroundValidator.RegisterValidation("required_if", func(fl validator.FieldLevel) bool {
		params := strings.Split(fl.Param(), " ")
		field2Name := params[0]
		field2ExpectedValue := params[1]
		field2Value := data[field2Name].(string)

		if field2Value != field2ExpectedValue {
			return true
		}
		// Required
		return fl.Field().String() != ""
	})
	rawErrors := r.playgroundValidator.ValidateMap(data, rules)
	if len(rawErrors) == 0 {
		return nil
	}

	err := helpers_error.ValidationErrors{}
	for fieldName, rawError := range rawErrors {
		playgroundValidationErrors, _ := rawError.(validator.ValidationErrors)
		for _, playgroundValidationError := range playgroundValidationErrors {
			item := helpers_error.ValidationErrorItem{
				Field: fieldName,
				Tag:   playgroundValidationError.Tag(),
				Param: playgroundValidationError.Param(),
			}
			err.ValidationErrors = append(err.ValidationErrors, item)
		}
	}
	return errors.Wrap(helpers_error.ValidationError, err.Error())
}
