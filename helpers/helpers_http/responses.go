package helpers_http

import (
	"encoding/json"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"net/http"
)

func SendResponseJson(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(response)
	helpers_error.PanicIfError(err)
}

func SendMultiResponse[D any](w http.ResponseWriter, statusCode int, data []D, total int) {
	SendResponseJson(w, statusCode, &MultiResponse[D]{
		Data:         data,
		RecordsTotal: total,
	})
}
func SendSingleResponse[D any](w http.ResponseWriter, statusCode int, data D) {
	SendResponseJson(w, statusCode, &SingleResponse[D]{
		Data: data,
	})
}
func SendErrorResponse(w http.ResponseWriter, statusCode int, err string) {
	SendResponseJson(w, statusCode, &ErrorResponse{
		Error: err,
	})
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}
type SingleResponse[D any] struct {
	Data D `json:"data"`
}
type MultiResponse[D any] struct {
	RecordsTotal int `json:"recordsTotal"`
	Data         []D `json:"data"`
}
