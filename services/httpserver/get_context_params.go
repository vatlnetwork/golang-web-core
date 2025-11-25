package httpserver

import "net/http"

func GetContextParams(req *http.Request) map[string]any {
	params := req.Context().Value(paramsKey)
	if params == nil {
		return map[string]any{}
	}

	paramsMap, ok := params.(map[string]any)
	if !ok {
		return map[string]any{}
	}

	return paramsMap
}
