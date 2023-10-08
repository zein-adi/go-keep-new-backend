package auth_handlers_restful

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"strconv"
	"time"
)

func NewRoleRestful(service auth_service_interfaces.IRoleServices) *RoleRestful {
	return &RoleRestful{service: service}
}

type RoleRestful struct {
	service auth_service_interfaces.IRoleServices
}

func (x *RoleRestful) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	request := auth_requests.NewGetRequest()
	q := r.URL.Query()

	if q.Has("skip") {
		request.Skip, _ = strconv.Atoi(q.Get("skip"))
	}
	if q.Has("take") {
		request.Take, _ = strconv.Atoi(q.Get("take"))
	}
	if q.Has("search") {
		request.Search = q.Get("search")
	}

	models := x.service.Get(ctx, request)
	count := x.service.Count(ctx, request)

	h.SendMultiResponse(w, http.StatusOK, models, count)
}
func (x *RoleRestful) Insert(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &auth_entities.Role{}
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
func (x *RoleRestful) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &auth_entities.Role{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}
	input.Id = p.ByName("roleId")

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
func (x *RoleRestful) DeleteById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	id := p.ByName("roleId")
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
