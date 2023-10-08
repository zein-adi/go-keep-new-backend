package helpers_error

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"strings"
)

/*
 * Multiple Errors
 */

var ValidationError = errors.New("validation error")

func NewValidationErrors(field, tag, param string) error {
	err := ValidationErrors{
		ValidationErrors: []ValidationErrorItem{
			{
				Field: field,
				Tag:   tag,
				Param: param,
			},
		},
	}
	return errors.Wrap(ValidationError, err.Error())
}

type ValidationErrors struct {
	ValidationErrors []ValidationErrorItem
}

func (r ValidationErrors) Error() string {
	errs := helpers.Map(r.ValidationErrors, func(d ValidationErrorItem) string {
		return d.Error()
	})
	return strings.Join(errs, "|")
}

/*
 * Single Errors
 */

type ValidationErrorItem struct {
	Field string
	Tag   string
	Param string
}

func (v ValidationErrorItem) Error() string {
	err := fmt.Sprintf("%s.%s", v.Field, v.Tag)
	if v.Param != "" {
		err += "." + v.Param
	}
	return err
}
