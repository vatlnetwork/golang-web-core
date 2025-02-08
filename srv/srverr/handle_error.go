package srverr

import (
	"golang-web-core/util"
	"net/http"
)

func Handle500(rw http.ResponseWriter, err error) {
	http.Error(rw, err.Error(), http.StatusInternalServerError)
	util.LogColor("red", "INTERNAL SERVER ERROR: %v", err.Error())
}

func Handle400(rw http.ResponseWriter, err error) {
	http.Error(rw, err.Error(), http.StatusBadRequest)
}
