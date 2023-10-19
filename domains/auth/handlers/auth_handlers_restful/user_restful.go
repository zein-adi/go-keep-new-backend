package auth_handlers_restful

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_services"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"strconv"
	"time"
)

func NewUserRestfulHandler(service auth_service_interfaces.IUserServices) *UserRestfulHandler {
	return &UserRestfulHandler{service: service}
}

type UserRestfulHandler struct {
	service auth_service_interfaces.IUserServices
}

func (x *UserRestfulHandler) Get(w http.ResponseWriter, r *http.Request) {
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
func (x *UserRestfulHandler) Insert(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &auth_requests.UserInputRequest{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}

	accessToken, _ := middlewares.GetAuthorizationToken(r)
	accessClaim, _ := middlewares.GetJwtClaims(accessToken)
	model, err := x.service.Insert(ctx, input, accessClaim.RoleIds)
	if err != nil {
		if errors.Is(err, helpers_error.ValidationError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else if errors.Is(err, auth_services.RoleAccessUnauthorizedError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}

	h.SendSingleResponse(w, http.StatusOK, model)
}
func (x *UserRestfulHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &auth_requests.UserUpdateRequest{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}
	vars := mux.Vars(r)
	input.Id = vars["userId"]

	accessToken, _ := middlewares.GetAuthorizationToken(r)
	accessClaim, _ := middlewares.GetJwtClaims(accessToken)
	model, err := x.service.Update(ctx, input, accessClaim.RoleIds)
	if err != nil {
		if errors.Is(err, helpers_error.EntryNotFoundError) {
			h.SendErrorResponse(w, http.StatusNotFound, "")
		} else if errors.Is(err, helpers_error.ValidationError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else if errors.Is(err, auth_services.RoleAccessUnauthorizedError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}

	h.SendSingleResponse(w, http.StatusOK, model)
}
func (x *UserRestfulHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &auth_requests.UserUpdatePasswordRequest{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}
	vars := mux.Vars(r)
	input.Id = vars["userId"]

	accessToken, _ := middlewares.GetAuthorizationToken(r)
	accessClaim, _ := middlewares.GetJwtClaims(accessToken)
	model, err := x.service.UpdatePassword(ctx, input, accessClaim.RoleIds)
	if err != nil {
		if errors.Is(err, helpers_error.EntryNotFoundError) {
			h.SendErrorResponse(w, http.StatusNotFound, "")
		} else if errors.Is(err, helpers_error.ValidationError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else if errors.Is(err, auth_services.RoleAccessUnauthorizedError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}

	h.SendSingleResponse(w, http.StatusOK, model)
}
func (x *UserRestfulHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["userId"]
	accessToken, _ := middlewares.GetAuthorizationToken(r)
	accessClaim, _ := middlewares.GetJwtClaims(accessToken)
	affected, err := x.service.DeleteById(ctx, id, accessClaim.RoleIds)
	if err != nil {
		if errors.Is(err, helpers_error.EntryNotFoundError) {
			h.SendErrorResponse(w, http.StatusNotFound, "")
		} else if errors.Is(err, auth_services.RoleAccessUnauthorizedError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}
	h.SendSingleResponse(w, http.StatusOK, affected)
}
