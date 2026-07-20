package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement basic auth
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				fmt.Errorf("authorization header is missing")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				next.ServeHTTP(w, r)
				fmt.Errorf("invalid authorization header format")
				return
			}

			// Decode the base64-encoded credentials
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				next.ServeHTTP(w, r)
				fmt.Errorf("failed to decode authorization header: %v", err)
				return
			}

			// Split the decoded string into username and password
			credentials := strings.SplitN(string(decoded), ":", 2)
			if len(credentials) != 2 {
				next.ServeHTTP(w, r)
				fmt.Errorf("invalid authorization header format")
				return
			}

			username := credentials[0]
			password := credentials[1]

			// Check if the username and password are valid
			if username != "admin" || password != "password" {
				next.ServeHTTP(w, r)
				fmt.Errorf("invalid username or password")
				return
			}

			// If we get here, the credentials are valid
			next.ServeHTTP(w, r)
		})
	}
}
