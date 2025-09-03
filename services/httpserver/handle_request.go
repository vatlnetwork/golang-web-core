package httpserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type paramsKeyType string

const paramsKey paramsKeyType = "params"

func handleOptions(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "*")
	rw.WriteHeader(http.StatusOK)
}

func (h *HttpServer) handleRequest(route Route) http.HandlerFunc {
	if route.Method == http.MethodOptions {
		return handleOptions
	}

	return func(rw http.ResponseWriter, req *http.Request) {
		requestId := uuid.NewString()
		req.Header.Set("X-Request-ID", requestId)

		remoteAddr := req.Header.Get("X-Forwarded-For")
		if remoteAddr == "" {
			remoteAddr = req.Header.Get("X-Real-IP")
		}
		if remoteAddr == "" {
			remoteAddr = req.RemoteAddr
		}
		h.logger.Infof("Started %v %v for %v", req.Method, req.URL.Path, remoteAddr)

		params, err := h.getParams(req)
		if err == nil {
			h.logger.Debugf("Params: %+v", params)
		}
		if params == nil {
			params = map[string]any{}
		}

		reqWithParams := req.WithContext(context.WithValue(req.Context(), paramsKey, params))
		route.Handler(rw, reqWithParams)

		h.logger.Infof("Completed %v %v for %v", req.Method, req.URL.Path, remoteAddr)
	}
}

func (h *HttpServer) getParams(req *http.Request, maxSize ...int64) (map[string]any, error) {
	size := int64(100)
	if len(maxSize) > 0 {
		size = maxSize[0]
	}

	if req.Method == http.MethodGet {
		queryValues := req.URL.Query()
		params := map[string]any{}
		for key, values := range queryValues {
			params[key] = values[0]
		}
		return params, nil
	}

	params, err := decodeRequestBodyToMap(req)
	if err != nil {
		params, err = decodeFormDataToMap(req, size)
		if err != nil {
			return nil, err
		}
	}

	return params, nil
}

func decodeRequestBodyToMap(req *http.Request) (map[string]any, error) {
	if req.Body == nil {
		return map[string]any{}, nil
	}

	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var decoded map[string]any
	err = json.Unmarshal(bytes, &decoded)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func decodeFormDataToMap(req *http.Request, maxSize ...int64) (map[string]any, error) {
	size := int64(100)
	if len(maxSize) > 0 {
		size = maxSize[0]
	}

	err := req.ParseMultipartForm(size << 20)
	if err != nil {
		return nil, err
	}

	decoded := map[string]any{}
	for key, value := range req.MultipartForm.Value {
		decoded[key] = value[0]
	}

	for key, value := range req.MultipartForm.File {
		decoded[key] = value[0]
	}

	return decoded, nil
}
