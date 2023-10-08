package middlewares

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"strings"
)

func Auth(writer http.ResponseWriter, request *http.Request, params httprouter.Params, routeName string) bool {
	tokenString, err := GetAuthorizationToken(request)
	if err != nil {
		helpers_http.SendErrorResponse(writer, http.StatusUnauthorized, err.Error())
		return false
	}
	_, err = GetJwtClaims(tokenString)
	if err != nil {
		helpers_http.SendErrorResponse(writer, http.StatusUnauthorized, err.Error())
		return false
	}
	return true
}

func AuthRefresh(writer http.ResponseWriter, request *http.Request, params httprouter.Params, routeName string) bool {
	tokenString, err := GetAuthorizationToken(request)
	if err != nil {
		helpers_http.SendErrorResponse(writer, http.StatusUnauthorized, err.Error())
		return false
	}
	_, err = GetRefreshJwtClaims(tokenString)
	if err != nil {
		helpers_http.SendErrorResponse(writer, http.StatusUnauthorized, err.Error())
		return false
	}
	return false
}

func GetAuthorizationToken(r *http.Request) (string, error) {
	bearerString := r.Header.Get("Authorization")
	bearer := strings.Split(bearerString, " ")
	if strings.ToLower(bearer[0]) != "bearer" || len(bearer) < 2 {
		return "", errors.New("bearer token expected")
	}

	return bearer[1], nil
}
func GetRefreshJwtClaims(tokenString string) (jwt.MapClaims, error) {
	// TODO get from env
	var refreshTokenSecret = "123456789"
	return getJwtClaims(tokenString, refreshTokenSecret)
}
func GetJwtClaims(tokenString string) (jwt.MapClaims, error) {
	// TODO get from env
	var accessTokenSecret = "123456789"
	return getJwtClaims(tokenString, accessTokenSecret)
}
func getJwtClaims(tokenString string, tokenSecret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid jwt signing method")
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid jwt token")
	}

	return claims, nil
}
