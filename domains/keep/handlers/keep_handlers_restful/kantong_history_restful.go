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
	request.Search = q.Get("search")
	request.Skip, err = strconv.Atoi(q.Get("skip"))
	helpers_error.PanicIfError(err)
	request.Take, err = strconv.Atoi(q.Get("take"))
	helpers_error.PanicIfError(err)

	vars := mux.Vars(r)
	kantongId := vars["kantongId"]
	models := x.service.Get(ctx, kantongId, request)
	h.SendMultiResponse(w, http.StatusOK, models, len(models))
}

func (x *KantongHistoryRestfulHandler) InsertAndUpdateSaldoKantong(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &keep_request.KantongHistoryInsertUpdate{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}

	vars := mux.Vars(r)
	input.KantongId = vars["kantongId"]
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
	input.KantongId = vars["kantongId"]
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
	kantongId := vars["kantongId"]
	id := vars["kantongHistoryId"]
	affected, err := x.service.DeleteById(ctx, kantongId, id)
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
