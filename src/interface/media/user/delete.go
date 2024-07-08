package media_user_int

import (
	"golang-web-core/src/application/srv/application/srverr"
	sessions_user_int "golang-web-core/src/interface/sessions/user"
	"net/http"

	"github.com/gorilla/mux"
)

func (c Controller) Delete(rw http.ResponseWriter, req *http.Request) {
	// find media file db entry
	id := mux.Vars(req)["id"]
	mediaFile, err := c.MediaRepository.Find(id)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// delete db entry
	currentUser, err := sessions_user_int.CurrentUser(req)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
	err = c.MediaRepository.Delete(mediaFile.Id, currentUser.Id)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}

	// delete file from storage
	err = c.OsMediaRepo.Delete(mediaFile, currentUser.Id)
	if err != nil {
		srverr.Raise(rw, req, err, http.StatusInternalServerError)
		return
	}
}
