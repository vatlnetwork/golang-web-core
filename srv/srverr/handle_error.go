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
	util.LogColor("yellow", "BAD REQUEST: %v", err.Error())
}

func Handle404(rw http.ResponseWriter, err error) {
	http.Error(rw, err.Error(), http.StatusNotFound)
	util.LogColor("yellow", "NOT FOUND: %v", err.Error())
}
