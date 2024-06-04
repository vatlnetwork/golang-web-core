package srverr

import (
	"log"
	"net/http"
)

func Raise(rw http.ResponseWriter, req *http.Request, err error, status int) {
	http.Error(rw, err.Error(), status)
	log.Println(err)
	log.Printf("%v %v responded with %v\n", req.Method, req.URL.Path, status)
}
