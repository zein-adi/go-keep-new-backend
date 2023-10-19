package keep_handlers_restful

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"time"
)

func NewKantongHistoryRestfulHandler(kantongHistoryService keep_service_interfaces.IKantongHistoryServices) *KantongHistoryRestfulHandler {
	return &KantongHistoryRestfulHandler{
		service: kantongHistoryService,
	}
}

type KantongHistoryRestfulHandler struct {
	service keep_service_interfaces.IKantongHistoryServices
}

func (x *KantongHistoryRestfulHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	models := x.service.Get(ctx)
	h.SendMultiResponse(w, http.StatusOK, models, len(models))
}

func (x *KantongHistoryRestfulHandler) Insert(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &keep_request.KantongHistoryInsertUpdate{}
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

func (x *KantongHistoryRestfulHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &keep_request.KantongHistoryInsertUpdate{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}

	vars := mux.Vars(r)
	input.Id = vars["kantongHistoryId"]
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

func (x *KantongHistoryRestfulHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["kantongHistoryId"]
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
