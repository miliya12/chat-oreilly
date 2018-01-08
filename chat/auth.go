package main

import (
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// Unauthorized
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error)
	} else {
		// Authorized
		h.next.ServeHTTP(w, r)
	}
}
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// check authorization result
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// set user data to cookie
	authCookieValue := objx.New(map[string]interface{}{
		"name":       user.Name,
		"avatar_url": user.AvatarURL,
	}).MustBase64()
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	})

	w.Header()["Location"] = []string{"/chat"}
	w.WriteHeader(http.StatusTemporaryRedirect)

}
