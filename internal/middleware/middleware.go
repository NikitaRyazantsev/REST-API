package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Logging(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func PanicRecovery(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Println(string(debug.Stack()))
			}
		}()
		h.ServeHTTP(w, r)
	}
}
