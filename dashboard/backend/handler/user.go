package handler

import (
	"net/http"
)

const (
	defaultAvatar = "https://gw.alipayobjects.com/zos/rmsportal/BiazfanxmamNRoxxVxka.png"
)

func (handler *APIHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("User")
	email := r.Header.Get("Email")

	u := User{
		Name:   user,
		Email:  email,
		Avatar: defaultAvatar,
	}

	responseJSON(u, w, http.StatusOK)
}
