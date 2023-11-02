package keep_handlers_restful

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
	"net/http"
	"strconv"
	"time"
)

func NewKantongRestfulHandler(posService keep_service_interfaces.IKantongServices) *KantongRestfulHandler {
	return &KantongRestfulHandler{
		service: posService,
	}
}

type KantongRestfulHandler struct {
	service keep_service_interfaces.IKantongServices
}

func (x *KantongRestfulHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	q := r.URL.Query()
	v := validator.New()
	err := v.ValidateMap(map[string]any{
		"skip": q.Get("skip"),
		"take": q.Get("take"),
	}, map[string]any{
		"skip": "omitempty,number",
		"take": "omitempty,number",
	})
	if err != nil {
		h.SendErrorResponse(w, 400, err.Error())
		return
	}

	request := helpers_requests.NewGet()
	request.Skip, err = strconv.Atoi(q.Get("skip"))
	helpers_error.PanicIfError(err)
	request.Take, err = strconv.Atoi(q.Get("take"))
	helpers_error.PanicIfError(err)
	request.Search = q.Get("search")

	err = v.ValidateStruct(request)
	if err != nil {
		h.SendErrorResponse(w, 400, err.Error())
		return
	}

	models := x.service.Get(ctx, request)
	h.SendMultiResponse(w, http.StatusOK, models, len(models))
}

func (x *KantongRestfulHandler) Insert(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &keep_request.KantongInsert{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}

	model, err := x.service.Insert(ctx, input)
	if err != nil {
		if errors.Is(err, helpers_error.ValidationError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}

	h.SendSingleResponse(w, http.StatusOK, model)
}

func (x *KantongRestfulHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &keep_request.KantongUpdate{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}

	vars := mux.Vars(r)
	input.Id = vars["kantongId"]
	model, err := x.service.Update(ctx, input)
	if err != nil {
		if errors.Is(err, helpers_error.EntryNotFoundError) {
			h.SendErrorResponse(w, http.StatusNotFound, "")
		} else if errors.Is(err, helpers_error.ValidationError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}

	h.SendSingleResponse(w, http.StatusOK, model)
}
func (x *KantongRestfulHandler) UpdateUrutan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	data := h.NewDefaultFormRequest([]*keep_request.KantongUpdateUrutanItem{})
	if !h.ReadRequest(w, r, data) {
		return
	}
	input := data.Data

	affected, err := x.service.UpdateUrutan(ctx, input)
	if err != nil {
		if errors.Is(err, helpers_error.ValidationError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}

	h.SendSingleResponse(w, http.StatusOK, affected)
}
func (x *KantongRestfulHandler) UpdateVisivility(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	data := h.NewDefaultFormRequest([]*keep_request.KantongUpdateVisibilityItem{})
	if !h.ReadRequest(w, r, data) {
		return
	}
	input := data.Data

	affected, err := x.service.UpdateVisibility(ctx, input)
	if err != nil {
		if errors.Is(err, helpers_error.ValidationError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}

	h.SendSingleResponse(w, http.StatusOK, affected)
}

func (x *KantongRestfulHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["kantongId"]
	affected, err := x.service.DeleteById(ctx, id)
	if err != nil {
		if errors.Is(err, helpers_error.EntryNotFoundError) {
			h.SendErrorResponse(w, http.StatusNotFound, "")
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}
	h.SendSingleResponse(w, http.StatusOK, affected)
}
func (x *KantongRestfulHandler) GetTrashed(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	q := r.URL.Query()
	v := validator.New()
	err := v.ValidateMap(map[string]any{
		"skip": q.Get("skip"),
		"take": q.Get("take"),
	}, map[string]any{
		"skip": "omitempty,number",
		"take": "omitempty,number",
	})
	if err != nil {
		h.SendErrorResponse(w, 400, err.Error())
		return
	}

	request := helpers_requests.NewGet()
	request.Skip, err = strconv.Atoi(q.Get("skip"))
	helpers_error.PanicIfError(err)
	request.Take, err = strconv.Atoi(q.Get("take"))
	helpers_error.PanicIfError(err)
	request.Search = q.Get("search")

	err = v.ValidateStruct(request)
	if err != nil {
		h.SendErrorResponse(w, 400, err.Error())
		return
	}

	models := x.service.GetTrashed(ctx, request)
	h.SendMultiResponse(w, http.StatusOK, models, len(models))
}
func (x *KantongRestfulHandler) RestoreTrashedById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["kantongId"]
	affected, err := x.service.RestoreTrashedById(ctx, id)
	if err != nil {
		if errors.Is(err, helpers_error.EntryNotFoundError) {
			h.SendErrorResponse(w, http.StatusNotFound, "")
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}
	h.SendSingleResponse(w, http.StatusOK, affected)
}
func (x *KantongRestfulHandler) DeleteTrashedById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["kantongId"]
	affected, err := x.service.DeleteTrashedById(ctx, id)
	if err != nil {
		if errors.Is(err, helpers_error.EntryNotFoundError) {
			h.SendErrorResponse(w, http.StatusNotFound, "")
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}
	h.SendSingleResponse(w, http.StatusOK, affected)
}
