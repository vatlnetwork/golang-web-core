package util

import (
	"net/http"
)

func GetParams(req *http.Request, maxSize ...int64) (map[string]any, error) {
	size := int64(100)
	if len(maxSize) > 0 {
		size = maxSize[0]
	}

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
		params, err = DecodeFormDataToMap(req, size)
		if err != nil {
			return nil, err
		}
	}

	return params, nil
}

func GetParamsFromContext(req *http.Request) map[string]any {
	return req.Context().Value("params").(map[string]any)
}
