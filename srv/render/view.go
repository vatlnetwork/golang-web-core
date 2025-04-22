package render

import (
	"fmt"
	"html/template"
	"inventory-app/srv/srverr"
	"net/http"
)

func RenderView(rw http.ResponseWriter, view string, data any) {
	viewPath := fmt.Sprintf("app/views/%v", view)
	template, err := template.ParseFiles(viewPath)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
	err = template.Execute(rw, data)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}
}
