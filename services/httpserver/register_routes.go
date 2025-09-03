package httpserver

import "fmt"

func (h *HttpServer) RegisterRoutes() {
	for _, route := range h.routes {
		h.mux.HandleFunc(fmt.Sprintf("%v %v", route.Method, route.Pattern), h.handleRequest(route))
	}
}
