package token

import (
	"math/rand"
	"net/http"
	"time"
	"unicode/utf8"
)

// Require is a simple HTTP Middleware that checks that the request comes with
// a cookie with the correct token. Otherwise the user is redirected to the
// given loginURL.
func Require(token string, loginURL string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("guideToken")
		if err != nil || c.Value != token {
			w.Header().Set("Location", loginURL)
			w.WriteHeader(307)
			return
		}
		f(w, r)
	}
}

// Generate produces a simple 6-character random-string that can be used as
// token.
func Generate() string {
	rand.Seed(time.Now().Unix())
	start, _ := utf8.DecodeRuneInString("a")
	end, _ := utf8.DecodeRuneInString("z")
	r := int(end) - int(start)
	result := ""
	for i := 0; i < 6; i++ {
		c := rand.Intn(r)
		result += string(rune(int(start) + c))
	}
	return result
}
