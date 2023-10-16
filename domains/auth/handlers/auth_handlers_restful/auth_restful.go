package auth_handlers_restful

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_responses"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"time"
)

func NewAuthRestful(service auth_service_interfaces.IAuthServices) *AuthRestful {
	return &AuthRestful{service: service}
}

type AuthRestful struct {
	service auth_service_interfaces.IAuthServices
}

func (x *AuthRestful) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, timeout := context.WithTimeout(context.Background(), time.Second*30)
	defer timeout()

	input := &auth_requests.LoginRequest{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}

	accessToken, refreshToken, err := x.service.Login(ctx, input.Username, input.Password, input.RememberMe)
	if err != nil {
		if errors.Is(err, helpers_error.ValidationError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}

	response := &auth_responses.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	h.SendSingleResponse(w, http.StatusOK, response)
}
func (x *AuthRestful) Refresh(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, timeout := context.WithTimeout(context.Background(), time.Second*30)
	defer timeout()

	refreshToken, _ := middlewares.GetAuthorizationToken(r)
	accessToken, updatedRefreshToken, err := x.service.Refresh(ctx, refreshToken)
	if err != nil {
		h.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	response := &auth_responses.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: updatedRefreshToken,
	}
	h.SendSingleResponse(w, http.StatusOK, response)
}
func (x *AuthRestful) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, timeout := context.WithTimeout(context.Background(), time.Second*30)
	defer timeout()

	refreshToken, _ := middlewares.GetAuthorizationToken(r)
	err := x.service.Logout(ctx, refreshToken)
	if err != nil {
		h.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.SendSingleResponse(w, http.StatusOK, "")
}
func (x *AuthRestful) Profile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, timeout := context.WithTimeout(context.Background(), time.Second*30)
	defer timeout()

	accessToken, _ := middlewares.GetAuthorizationToken(r)
	response, err := x.service.Profile(ctx, accessToken)
	if err != nil {
		h.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	h.SendSingleResponse(w, http.StatusOK, response)
}
