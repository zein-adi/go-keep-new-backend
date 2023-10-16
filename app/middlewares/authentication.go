package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"strings"
	"time"
)

var (
	JwtTokenError = errors.New("jwt token error")
)

// TODO get from env
var (
	refreshTokenSecret = "123456789"
	accessTokenSecret  = "123456789"
)

func AuthHandle(writer http.ResponseWriter, request *http.Request, _ httprouter.Params, _ string) bool {
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

func AuthRefreshHandle(writer http.ResponseWriter, request *http.Request, _ httprouter.Params, _ string) bool {
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
	return true
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
	return getJwtClaims(tokenString, refreshTokenSecret)
}
func GetJwtClaims(tokenString string) (jwt.MapClaims, error) {
	return getJwtClaims(tokenString, accessTokenSecret)
}
func getJwtClaims(tokenString string, tokenSecret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.Wrap(JwtTokenError, "invalid jwt signing method")
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.Wrap(JwtTokenError, "invalid jwt token")
	}

	return claims, nil
}
func IssueNewAccessToken(username, nama string, roleIds []string, iat, exp time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     username,
		"name":    nama,
		"roleIds": roleIds,
		"iat":     iat.Unix(),
		"exp":     exp.Unix(),
	})
	tokenString, err := token.SignedString([]byte(accessTokenSecret))
	helpers_error.PanicIfError(err)
	return tokenString
}
func IssueNewRefreshToken(username string, rememberMe bool, iat, exp time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      username,
		"remember": rememberMe,
		"iat":      iat.Unix(),
		"exp":      exp.Unix(),
	})
	tokenString, err := token.SignedString([]byte(refreshTokenSecret))
	helpers_error.PanicIfError(err)
	return tokenString
}
