package util

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeRequestBody(req *http.Request, decodeObject interface{}) error {
	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, decodeObject)
	if err != nil {
		return err
	}

	return nil
}
