package util

import (
	"net/http"
)

func GetParams(req *http.Request) (map[string]any, error) {
	if req.Method == http.MethodGet {
		queryValues := req.URL.Query()
		params := make(map[string]any)
		for key, value := range queryValues {
			params[key] = value[0]
		}
		return params, nil
	}

	params, err := DecodeRequestBodyToMap(req)
	if err != nil {
		params, err = DecodeFormDataToMap(req)
		if err != nil {
			return nil, err
		}
	}

	return params, nil
}
