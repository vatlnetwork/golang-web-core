// this is an example controller

package controllers

import (
	"golang-web-core/srv/cfg"
	"golang-web-core/srv/render"
	"net/http"
	"reflect"
)

type TestController struct {
	cfg.Config
}

func NewTestController(c cfg.Config) TestController {
	return TestController{
		Config: c,
	}
}

func (c TestController) Name() string {
	return reflect.TypeOf(c).Name()
}

func (c TestController) BeforeAction(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}
}

func (c TestController) TestMethod(rw http.ResponseWriter, req *http.Request) {
	render.RenderView(rw, "test/test_method.go.tmpl", "If you see this message, it means the test method worked.")
}

func (c TestController) TestMemberMethod(rw http.ResponseWriter, req *http.Request) {
	member_var := req.PathValue("member_var")

	render.RenderView(rw, "test/test_member_method.go.tmpl", member_var)
}
