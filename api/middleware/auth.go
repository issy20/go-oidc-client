package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/issy20/go-oidc-client/util"
)

type ContextKey string

type AuthorizationHeader struct {
	AccessToken string
	ClientID    string
}

var AutorizationContextKey = ContextKey("authorization")

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func AuthMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("sub")
		if err != nil {
			log.Print("AuthMiddleware ", err)
		}
		log.Print(cookie)

		tokenStr := ExtractToken(r)
		secrets, err := util.ReadJson()

		authorizationHeader := &AuthorizationHeader{
			AccessToken: tokenStr,
			ClientID:    secrets.ClientId,
		}

		if err != nil {
			log.Fatal(err)
		}

		ctx := context.WithValue(r.Context(), AutorizationContextKey, authorizationHeader)
		r = r.WithContext(ctx)
		f.ServeHTTP(w, r)
	})
}
