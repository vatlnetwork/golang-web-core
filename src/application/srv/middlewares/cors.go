package middlewares

import "net/http"

func EnableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(rw, req)
	})
}
