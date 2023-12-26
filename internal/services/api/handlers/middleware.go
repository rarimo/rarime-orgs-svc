package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strings"
)

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				ape.Render(w, problems.Unauthorized())
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")
			// TODO: implement auth and user did extraction and it's context injection

			next.ServeHTTP(w, r)
		})
	}
}
