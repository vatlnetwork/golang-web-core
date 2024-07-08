package media_user_int

import (
	"encoding/json"
	"golang-web-core/src/application/srv/application/srverr"
	"net/http"

	"github.com/gorilla/mux"
)

func (c Controller) UserMedia(rw http.ResponseWriter, req *http.Request) {
	// get current user's media
	id := mux.Vars(req)["id"]
	media, err := c.MediaRepository.FindByUser(id)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// return current user's media
	res, err := json.Marshal(media)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	rw.Write(res)
}
