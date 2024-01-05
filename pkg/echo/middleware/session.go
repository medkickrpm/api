package middleware

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

func Setup() {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 1, // 1 day
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode, // This sets the SameSite attribute
		Secure:   false,                   // Only set this if you're using HTTPS!
	}
}
