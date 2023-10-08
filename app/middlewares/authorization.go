package middlewares

import (
	"github.com/julienschmidt/httprouter"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"strings"
)

func Acl(writer http.ResponseWriter, request *http.Request, params httprouter.Params, routeName string) bool {
	tokenString, _ := GetAuthorizationToken(request)
	claims, _ := GetJwtClaims(tokenString)

	roles := strings.Split(claims["roles"].(string), ",")
	permissions := getPermissionsByRole(roles...)

	_, err := helpers.FindIndex(permissions, func(s string) bool {
		return s == routeName
	})
	if err != nil {
		helpers_http.SendErrorResponse(writer, http.StatusForbidden, "")
		return false
	}
	return true
}

func getPermissionsByRole(roles ...string) []string {
	// TODO implement this
	return nil
}
