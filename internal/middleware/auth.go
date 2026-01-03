package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var JWTSecret = []byte("new-secret-key-change-this-in-production-2026") // TODO: Move to config

// Session durations
const (
	SessionDurationShort = 12 * time.Hour      // Sin "Recordarme"
	SessionDurationLong  = 30 * 24 * time.Hour // Con "Recordarme"
)

type Claims struct {
	Username   string `json:"username"`
	Role       string `json:"role"`
	RememberMe bool   `json:"remember_me"`
	jwt.RegisteredClaims
}

func AuthMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					logger.WithField("panic", r).Error("Panic en AuthMiddleware")
				}
			}()

			// Obtener token de cookie o header
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				cookie, err := r.Cookie("token")
				if err != nil {
					logger.Debug("No token encontrado, redirigiendo a login")
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}
				tokenString = cookie.Value
			} else if strings.HasPrefix(tokenString, "Bearer ") {
				tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			}

			// Validar token JWT
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return JWTSecret, nil
			})

			if err != nil || !token.Valid {
				logger.WithError(err).Debug("Token inválido")
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Extraer claims y continuar
			if claims, ok := token.Claims.(*Claims); ok && claims.Username != "" {
				ctx := context.WithValue(r.Context(), "username", claims.Username)
				ctx = context.WithValue(ctx, "role", claims.Role)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				logger.Debug("Claims inválidos en token")
				http.Redirect(w, r, "/login", http.StatusFound)
			}
		})
	}
}
