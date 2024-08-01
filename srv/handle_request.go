package srv

import (
	"golang-web-core/app/routes"
	"net/http"
)

func HandleRequest(route routes.Route) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
	}
}

func HandleOptions(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "*")
	rw.WriteHeader(http.StatusOK)
}
