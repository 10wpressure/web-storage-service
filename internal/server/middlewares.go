package server

import (
	"context"
	"net/http"
	"web-storage-service/pkg"
)

type key int

const UserIDKey key = 0

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := pkg.ExtractBearerToken(r)
		if err != nil {
			UnauthorizedError(w, "invalid authorization token")
			return
		}

		userID, err := s.DeactivateExpiredSessionsAndReturnUserID(token)
		if err != nil {
			UnauthorizedError(w, "invalid authorization token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
