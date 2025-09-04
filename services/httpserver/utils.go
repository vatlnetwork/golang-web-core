package httpserver

import (
	"encoding/json"
	"net/http"
)

func DecodeContextParams(req *http.Request, object any) error {
	params := GetParamsFromContext(req)

	paramsJson, err := json.Marshal(params)
	if err != nil {
		return err
	}

	err = json.Unmarshal(paramsJson, object)
	if err != nil {
		return err
	}

	return nil
}

func GetParamsFromContext(req *http.Request) map[string]any {
	return req.Context().Value(paramsKey).(map[string]any)
}
