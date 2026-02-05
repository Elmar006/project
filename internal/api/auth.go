package api

import (
	"net/http"
)

var todoPassword string

func PasswordCheck(password string) {
	todoPassword = password
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(todoPassword) > 0 {
			var jwt string

			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			var valid bool

			if jwt != "" {
				valid = validateToken(jwt, todoPassword)
			}
			if !valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
