package middlewares

import (
	"log"
	"net/http"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("[%v] %v %v", req.RemoteAddr, req.Method, req.URL.Path)
		next.ServeHTTP(rw, req)
	})
}
