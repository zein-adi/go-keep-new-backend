package basic_handlers_restful

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"strconv"
	"time"
)

func NewChangelogRestfulHandler(changelogService basic_service_interfaces.IChangelogServices) *ChangelogRestfulHandler {
	return &ChangelogRestfulHandler{
		service: changelogService,
	}
}

type ChangelogRestfulHandler struct {
	service basic_service_interfaces.IChangelogServices
}

func (x *ChangelogRestfulHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	q := r.URL.Query()
	take := 10
	skip := 0
	if q.Has("take") {
		t, err := strconv.Atoi(q.Get("take"))
		if err != nil {
			h.SendSingleResponse(w, 400, helpers_error.NewValidationErrors("take", "type", "integer"))
			return
		}
		take = t
	}
	if q.Has("skip") {
		t, err := strconv.Atoi(q.Get("skip"))
		if err != nil {
			h.SendSingleResponse(w, 400, helpers_error.NewValidationErrors("skip", "type", "integer"))
			return
		}
		skip = t
	}

	models := x.service.Get(ctx, skip, take)
	h.SendMultiResponse(w, http.StatusOK, models, len(models))
}
func (x *ChangelogRestfulHandler) Insert(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &basic_entities.Changelog{}
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
func (x *ChangelogRestfulHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &basic_entities.Changelog{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}

	vars := mux.Vars(r)
	input.Id = vars["changelogId"]
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
func (x *ChangelogRestfulHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["changelogId"]
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
