package hm

import (
	"net/http"
)

const (
	ErrFDuplicateCookie = "duplicate cookie: %s"
)

// SetCookie sets a cookie.
// Unfortunately, Go language is so ugly that it does not allow to check
// whether the specified cookie is already set or not.
func SetCookie(rw http.ResponseWriter, cookie *http.Cookie) {
	http.SetCookie(rw, cookie)
}

// UnsetCookie unsets a cookie.
func UnsetCookie(rw http.ResponseWriter, cookie *http.Cookie) {
	cookie.Value = ""

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie#max-agenumber
	cookie.MaxAge = 0

	http.SetCookie(rw, cookie)
}

//----------------------------------------------------------------------------//
// Q:	What do you like in Go language ?
//
// A:	Go language is a good example showing what a programming language
//		should NOT be. Moreover, by looking at how people use this language one
//		can learn their character and attitude to this world.
//----------------------------------------------------------------------------//
