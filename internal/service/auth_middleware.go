package service

import (
	"net/http"
	"strings"

	"github.com/Stern-Ritter/go_task_manager/internal/utils"
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rootPassword := s.Config.RootPassword

		if len(strings.TrimSpace(rootPassword)) > 0 {
			var jwt string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			err = utils.CompareHash(rootPassword, jwt)

			if err != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
