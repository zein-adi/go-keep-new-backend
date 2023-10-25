package basic_handlers_restful

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
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
	v := validator.New()
	err := v.ValidateMap(map[string]interface{}{
		"skip": q.Get("skip"),
		"take": q.Get("take"),
	}, map[string]interface{}{
		"skip": "omitempty,number",
		"take": "omitempty,number",
	})
	if err != nil {
		h.SendSingleResponse(w, 400, err.Error())
		return
	}

	request := helpers_requests.NewGet()
	request.Search = q.Get("search")
	if q.Has("skip") {
		request.Skip, _ = strconv.Atoi(q.Get("skip"))
	}
	if q.Has("take") {
		request.Take, _ = strconv.Atoi(q.Get("take"))
	} else {
		request.Take = 10
	}
	err = v.ValidateStruct(request)
	if err != nil {
		h.SendSingleResponse(w, 400, err.Error())
		return
	}

	models := x.service.Get(ctx, request)
	count := x.service.Count(ctx, request)
	h.SendMultiResponse(w, http.StatusOK, models, count)
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
	affected, err := x.service.Update(ctx, input)
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

	h.SendSingleResponse(w, http.StatusOK, affected)
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
