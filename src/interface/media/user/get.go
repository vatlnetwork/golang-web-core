package media_user_int

import (
	"golang-web-core/src/application/srv/application/srverr"
	"net/http"

	"github.com/gorilla/mux"
)

func (c Controller) Get(rw http.ResponseWriter, req *http.Request) {
	// find media file in database
	id := mux.Vars(req)["id"]
	mediaFile, err := c.MediaRepository.Find(id)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// read media file from storage
	bytes, err := c.OsMediaRepo.Load(mediaFile)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// write media bytes to response writer
	rw.Write(bytes)
}
