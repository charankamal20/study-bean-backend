package initializers

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

const (
	key = "5BwtdtuMMnMuGgYxpWjLnErRJyk="
	MaxAge = 86400 * 30
	IsProd = false
)

func NewAuth() {
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store 
}