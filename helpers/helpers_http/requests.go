package helpers_http

import (
	"encoding/json"
	"fmt"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

func ReadRequest(w http.ResponseWriter, r *http.Request, model any) bool {
	contentType := r.Header.Get("Content-Type")

	var err error
	if contentType == "application/json" {
		err = ReadRequestJson(r, model)
	} else {
		err = ReadRequestPost(r, model)
	}

	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}
func ReadRequestJson(r *http.Request, model any) error {
	body, err := io.ReadAll(r.Body)
	helpers_error.PanicIfError(err)
	err = json.Unmarshal(body, model)
	val, ok := err.(*json.UnmarshalTypeError)
	if ok && val.Field != "" && val.Type.Name() != "" && val.Value != "" {
		return fmt.Errorf("invalid request: %s expected %s actual %s",
			val.Field,
			val.Type.Name(),
			val.Value,
		)
	}
	return err
}
func ReadRequestPost(r *http.Request, model any) error {
	rv := reflect.ValueOf(model).Elem()
	rt := reflect.TypeOf(model).Elem()
	for i := 0; i < rt.NumField(); i++ {
		jsonField := rt.Field(i).Tag.Get("json")
		value := r.PostFormValue(jsonField)
		if r.PostForm.Has(jsonField) {
			structField := rt.Field(i).Name
			f := rv.FieldByName(structField)
			if f.IsValid() && f.CanSet() {
				switch f.Kind() {
				case reflect.String:
					f.SetString(value)
				case reflect.Int:
					x, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						return err
					}
					f.SetInt(x)
				}
			}
		}
	}
	return nil
}

func NewDefaultFormRequest[D any](data D) *DefaultFormRequest[D] {
	return &DefaultFormRequest[D]{
		Data: data,
	}
}

type DefaultFormRequest[D any] struct {
	Data D `json:"data"`
}
