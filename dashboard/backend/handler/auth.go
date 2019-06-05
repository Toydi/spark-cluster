package handler

import (
	"context"
	"fmt"
	"net/http"
)

const (
	AuthProviderOAuth = "oauth"
	AuthProviderLocal = "local"
)

func (handler *APIHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	if errMsg := r.FormValue("error"); errMsg != "" {
		http.Error(w, fmt.Sprintf("%v: %v", errMsg, r.FormValue("error_description")), http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	if len(code) == 0 {
		http.Error(w, "No code in request", http.StatusBadRequest)
		return
	}

	// exchange token by grant code
	token, err := handler.oauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authenticate failed: %v", err), http.StatusBadRequest)
		return
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id token information in token", http.StatusBadRequest)
	}
	w.Header().Set("Location", fmt.Sprintf("/?access_token=%s", rawIDToken))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
