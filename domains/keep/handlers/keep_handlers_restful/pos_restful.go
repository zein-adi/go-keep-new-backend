package keep_handlers_restful

import (
	"context"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"time"
)

func NewPosRestfulHandler(posService keep_service_interfaces.IPosServices) *PosRestfulHandler {
	return &PosRestfulHandler{
		service: posService,
	}
}

type PosRestfulHandler struct {
	service keep_service_interfaces.IPosServices
}

func (x *PosRestfulHandler) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	request := keep_request.NewPosGetRequest()
	q := r.URL.Query()
	if q.Has("isLeafOnly") {
		request.IsLeafOnly = true
	}

	models := x.service.Get(ctx, request)
	h.SendMultiResponse(w, http.StatusOK, models, len(models))
}
func (x *PosRestfulHandler) Insert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &keep_request.PosInputUpdateRequest{}
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
func (x *PosRestfulHandler) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &keep_request.PosInputUpdateRequest{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}
	input.Id = p.ByName("posId")

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
func (x *PosRestfulHandler) DeleteById(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	id := p.ByName("posId")
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
func (x *PosRestfulHandler) GetTrashed(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	models := x.service.GetTrashed(ctx)
	h.SendMultiResponse(w, http.StatusOK, models, len(models))
}
func (x *PosRestfulHandler) RestoreTrashedById(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	id := p.ByName("posId")
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
func (x *PosRestfulHandler) DeleteTrashedById(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	id := p.ByName("posId")
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
