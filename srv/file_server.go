package srv

import (
	"log"
	"net/http"
)

type FileServer struct {
	http.Handler
	Prefix string
}

func (s FileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Serving %v%v to %v\n", s.Prefix, req.URL.Path, req.RemoteAddr)
	s.Handler.ServeHTTP(rw, req)
}
