package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"regexp"
	"strings"
)

// New PlaygroundValidator Wrapper
func New() *Validator {
	v := &Validator{
		playgroundValidator: validator.New(),
	}
	v.registerValidPassword()
	return v
}

type Validator struct {
	playgroundValidator *validator.Validate
}

func (r *Validator) ValidateMap(data map[string]interface{}, rules map[string]interface{}) error {
	r.validateMapRegisterRequiredIf(data)

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

func (r *Validator) ValidateStruct(s interface{}) error {
	rawError := r.playgroundValidator.Struct(s)
	if rawError == nil {
		return rawError
	}

	err := helpers_error.ValidationErrors{}
	playgroundValidationErrors, _ := rawError.(validator.ValidationErrors)
	for _, playgroundValidationError := range playgroundValidationErrors {

		item := helpers_error.ValidationErrorItem{
			Field: strcase.ToSnake(playgroundValidationError.Field()),
			Tag:   playgroundValidationError.Tag(),
			Param: playgroundValidationError.Param(),
		}
		err.ValidationErrors = append(err.ValidationErrors, item)
	}
	return errors.Wrap(helpers_error.ValidationError, err.Error())
}

func (r *Validator) registerValidPassword() {
	err := r.playgroundValidator.RegisterValidation("valid_password", func(fl validator.FieldLevel) bool {
		currentValue := fl.Field().String()

		regex, _ := regexp.Compile(`[\d]`)
		containsNumber := regex.FindString(currentValue) != ""

		regex, _ = regexp.Compile(`[a-z]`)
		containsLowerCase := regex.FindString(currentValue) != ""

		regex, _ = regexp.Compile(`[A-Z]`)
		containsUpperCase := regex.FindString(currentValue) != ""

		return containsNumber && containsLowerCase && containsUpperCase
	})
	helpers_error.PanicIfError(err)
}
func (r *Validator) validateMapRegisterRequiredIf(data map[string]interface{}) {
	err := r.playgroundValidator.RegisterValidation("required_if", func(fl validator.FieldLevel) bool {
		params := strings.Split(fl.Param(), " ")
		field2Name := params[0]
		field2ExpectedValue := params[1]
		field2Value := data[field2Name].(string)

		if field2Value != field2ExpectedValue {
			return true
		}
		// Required
		currentValue := fl.Field().String()
		return currentValue != ""
	})
	helpers_error.PanicIfError(err)
}
