package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vickon16/third-backend-tutorial/cmd/sqlc"
	"github.com/vickon16/third-backend-tutorial/cmd/types"
	"github.com/vickon16/third-backend-tutorial/cmd/utils"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}

func JWT(db *sqlc.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("no token provided"))
				return
			}

			// remove the Bearer and leave the auth token string
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &types.JWTClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(utils.Configs.JWT_SECRET), nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token signature"))
					return
				}
				utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid token"))
				return
			}

			if !token.Valid {
				utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid token"))
				return
			}

			// check if user is in the database
			_, err = db.GetUserById(r.Context(), claims.UserId)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("UnAuthorized Access"))
				return
			}

			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
