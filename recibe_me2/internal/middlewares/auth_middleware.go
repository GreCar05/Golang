package middlewares

import (
	"net/http"
	"recibe_me/configs/constants"
	"strings"

	"recibe_me/internal/helpers"
	"recibe_me/pkg/crypto"
)

// Authenticate is a middleware for validate and authenticate a user, takes a handler and return another
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		token := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
		err := crypto.ValidateToken([]byte(token))

		if err == nil {
			next(responseWriter, request)
		} else {
			helpers.Response(responseWriter, http.StatusUnauthorized, constants.ERR_INVALID_TOKEN, err, nil)
		}
	}
}
