package middlewares

import (
	"chat/internal/domain/service"
	"log"
	"net/http"
)

func NewAuthMiddleware(
	auth *service.AuthService,
) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authToken, err := r.Cookie("chat_uid")
			if err != nil {
				log.Printf("failed to auth: %w\n", err)
			}
			if authToken.Value == "" {
				log.Println(http.StatusUnauthorized)
				log.Println("Please register and log in")
				return
			}

			isValid := auth.IsTokenValid(authToken.Value)
			if !isValid {
				log.Println(http.StatusUnauthorized)
				log.Println("Invalid token: Please register and log in")
				return
			}

			h.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
