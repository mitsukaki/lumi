package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mitsukaki/lumi/internal/db"
	"github.com/mitsukaki/lumi/models"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	Id string `json:"id"`

	jwt.RegisteredClaims
}

type AuthHandler struct {
	logger *zap.Logger

	userTable *db.CouchTable
}

func CreateAuthHandler(logger *zap.Logger, userTable *db.CouchTable) *AuthHandler {
	return &AuthHandler{
		logger:    logger,
		userTable: userTable,
	}
}

func (auth *AuthHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT from the "authorization" cookie
		jwtCookie, err := r.Cookie("authorization")
		if err != nil {
			auth.logger.Error("no auth cookie set")
			ctx := context.WithValue(r.Context(), "authenticated", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Parse the JWT
		token, err := jwt.ParseWithClaims(jwtCookie.Value, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(viper.GetString("api.hmac_secret")), nil
		})

		if err != nil {
			auth.logger.Error("failed to parse token", zap.Error(err))
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// validate the essential claims
		if !token.Valid {
			auth.logger.Error("token is invalid")
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Get the user from the database
		claims := token.Claims.(*AuthClaims)
		doc, err := auth.userTable.Get(claims.Id)
		if err != nil {
			auth.logger.Error("Failed to get user from database", zap.Error(err))
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Unmarshal the user from the database
		var dbUser models.DBUser
		err = json.Unmarshal(doc, &dbUser)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			auth.logger.Info("failed to unmarshal db user", zap.Error(err))
			return
		}

		// check the response from the database
		if dbUser.Error != "" {
			http.Error(w, http.StatusText(404), 404)
			auth.logger.Info("failed to get user", zap.Error(err))
			return
		}

		// Add the user to the context
		ctx := context.WithValue(r.Context(), "user", &dbUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
