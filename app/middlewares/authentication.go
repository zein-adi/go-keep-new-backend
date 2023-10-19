package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"strings"
	"time"
)

var (
	JwtTokenError = errors.New("jwt token error")
)

var (
	refreshTokenSecret = viper.GetString("AUTH_ACCESS_TOKEN_SECRET")
	accessTokenSecret  = viper.GetString("AUTH_REFRESH_TOKEN_SECRET")
)

func AuthHandle(writer http.ResponseWriter, request *http.Request, _ string) bool {
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

func AuthRefreshHandle(writer http.ResponseWriter, request *http.Request, _ string) bool {
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
func GetRefreshJwtClaims(tokenString string) (*RefreshTokenClaims, error) {
	claims, err := getJwtClaims(tokenString, refreshTokenSecret)
	if err != nil {
		return nil, err
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}
	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}
	iat, err := claims.GetIssuedAt()
	if err != nil {
		return nil, err
	}
	rememberAny, ok := claims["remember"]
	if !ok {
		return nil, errors.Wrap(JwtTokenError, "jwt token must contains remember field")
	}
	remember := rememberAny.(bool)

	model := &RefreshTokenClaims{
		Sub:      subject,
		Remember: remember,
		Iat:      iat.Unix(),
		Exp:      expirationTime.Unix(),
	}
	return model, nil
}
func GetJwtClaims(tokenString string) (*AccessTokenClaims, error) {
	claims, err := getJwtClaims(tokenString, accessTokenSecret)
	if err != nil {
		return nil, err
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}
	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}
	iat, err := claims.GetIssuedAt()
	if err != nil {
		return nil, err
	}
	nameInf, ok := claims["name"]
	if !ok {
		return nil, errors.Wrap(JwtTokenError, "jwt token must contains name field")
	}
	roleIdsInf, ok := claims["roleIds"]
	if !ok {
		return nil, errors.Wrap(JwtTokenError, "jwt token must contains roleIds field")
	}
	nameString := nameInf.(string)
	roleIdsAny := roleIdsInf.([]any)
	roleIds := helpers.Map(roleIdsAny, func(d any) string {
		return d.(string)
	})

	model := &AccessTokenClaims{
		Sub:     subject,
		Name:    nameString,
		RoleIds: roleIds,
		Iat:     iat.Unix(),
		Exp:     expirationTime.Unix(),
	}
	return model, nil
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
