package srv

import (
	"log"
	"net/http"
)

type FileServer struct {
	http.Handler
}

func (s FileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Serving /public/%v to %v\n", req.URL.Path, req.RemoteAddr)
	s.Handler.ServeHTTP(rw, req)
}
