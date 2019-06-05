package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (handler *APIHandler) AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// debug hack
	debug := r.Header.Get("Debug")
	debugUser := r.Header.Get("Debug-User")
	if len(debug) != 0 && len(debugUser) != 0 {
		r.Header.Add("User", debugUser)
		next(w, r)
		return
	}

	rawToken := r.Header.Get("Authorization")
	if len(rawToken) == 0 {
		http.Error(w, "Authorization required.", http.StatusForbidden)
		return
	}

	token, err := handler.verifier.Verify(context.Background(), rawToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Token error: %v", err), http.StatusForbidden)
		return
	}

	var claims struct {
		Email  string `json:"email"`
		User   string `json:"name"`
		Expiry int64  `json:"exp"`
	}
	if err := token.Claims(&claims); err != nil {
		http.Error(w, fmt.Sprintf("Token error: %v", err), http.StatusForbidden)
		return
	}

	if claims.Expiry < time.Now().Unix() {
		http.Error(w, "Token expired.", http.StatusForbidden)
		return
	}

	r.Header.Add("User", claims.User)
	r.Header.Add("Email", claims.Email)

	next(w, r)
}
