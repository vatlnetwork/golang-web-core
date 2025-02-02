package srverr

import "net/http"

func Handle500(rw http.ResponseWriter, err error) {
	http.Error(rw, err.Error(), http.StatusInternalServerError)
}

func Handle400(rw http.ResponseWriter, err error) {
	http.Error(rw, err.Error(), http.StatusBadRequest)
}
