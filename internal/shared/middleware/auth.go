package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
	"time-management/internal/shared/util"
	"time-management/internal/user/application/command"
	"time-management/internal/user/domain"
	userHttp "time-management/internal/user/interface/http"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(userHttp.CookieAuthName)
		if err != nil {
			_ = util.WriteJson(w, http.StatusUnauthorized, util.ApiError{Error: ErrNoValidToken.Error()})
			return
		}

		// Validate the token from the cookie
		claims, err := validateToken(cookie.Value)
		if err != nil {
			_ = util.WriteJson(w, http.StatusUnauthorized, util.ApiError{Error: ErrInvalidToken.Error()})
			return
		}

		// Add the extracted user to the context
		user := &domain.User{
			Id:   claims[command.JwtId].(string),
			Role: claims[command.JwtRole].(string),
		}

		// Add user to the context
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure token uses HMAC for signing
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnExpectedSigningMethod
		}

		// Return the secret key for validating the token's signature
		return jwtSecretKey, nil
	})

	// Check if there was an error in parsing the token
	if err != nil {
		fmt.Println("Token parsing error:", err)
		return nil, err
	}

	// Check if the token is valid and cast claims to jwt.MapClaims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Optionally, validate the claims such as expiration
		err := validateClaims(claims)
		if err != nil {
			fmt.Println("Token claims validation error:", err)
			return nil, err
		}

		// If everything is fine, return the claims
		return claims, nil
	}

	// If the token is invalid
	return nil, ErrInvalidToken
}

// validateClaims checks token claims such as expiration
func validateClaims(claims jwt.MapClaims) error {
	// Check if the token is expired
	if exp, ok := claims[command.JwtExp].(float64); ok {
		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			return ErrTokenExpired
		}
	}

	return nil
}
