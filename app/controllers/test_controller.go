// this is an example controller

package controllers

import (
	"encoding/json"
	"fmt"
	"golang-web-core/app/models"
	"golang-web-core/srv/cfg"
	"golang-web-core/srv/render"
	"golang-web-core/srv/route"
	"golang-web-core/srv/srverr"
	"golang-web-core/util"
	"net/http"
	"reflect"
)

type TestController struct {
	cfg.Config
}

// this verifies that TestController fully implements Controller
var TestControllerVerifier Controller = TestController{}

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
		// see the comments in application controller to understand more about BeforeAction
		handler(rw, req)
	}
}

func (c TestController) Routes() []route.Route {
	return []route.Route{
		{
			Pattern:        "/test_route",
			Method:         http.MethodGet,
			Handler:        c.TestMethod,
			ControllerName: c.Name(),
		},
		{
			Pattern:        "/test_member_route/{member_var}/test",
			Method:         http.MethodGet,
			Handler:        c.TestMemberMethod,
			ControllerName: c.Name(),
		},
		{
			Pattern:        "/test/create",
			Method:         http.MethodPost,
			Handler:        c.TestCreateMethod,
			ControllerName: c.Name(),
		},
		{
			Pattern:        "/test/read",
			Method:         http.MethodGet,
			Handler:        c.TestReadMethod,
			ControllerName: c.Name(),
		},
		{
			Pattern:        "/test/update",
			Method:         http.MethodPatch,
			Handler:        c.TestUpdateMethod,
			ControllerName: c.Name(),
		},
		{
			Pattern:        "/test/update",
			Method:         http.MethodPut,
			Handler:        c.TestUpdateMethod,
			ControllerName: c.Name(),
		},
		{
			Pattern:        "/test/{id}/delete",
			Method:         http.MethodDelete,
			Handler:        c.TestDestroyMethod,
			ControllerName: c.Name(),
		},
		{
			Pattern:        "/crud_test",
			Method:         http.MethodGet,
			Handler:        c.CRUDTest,
			ControllerName: c.Name(),
		},
	}
}

func (c TestController) TestMethod(rw http.ResponseWriter, req *http.Request) {
	render.RenderView(rw, "test/test_method.go.tmpl", "If you see this message, it means the test method worked.")
}

func (c TestController) TestMemberMethod(rw http.ResponseWriter, req *http.Request) {
	member_var := req.PathValue("member_var")

	render.RenderView(rw, "test/test_member_method.go.tmpl", member_var)
}

// these routes are for testing the CRUD functionality of the golang web core
// this section of the controller interacts with the models section of the web core

func (c TestController) TestCreateMethod(rw http.ResponseWriter, req *http.Request) {
	var inputObject models.TestObject
	err := util.DecodeRequestBody(req, &inputObject)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	object := models.NewTestObject(inputObject.Number, inputObject.Boolean)

	testModel := models.NewTestModel(&c.Config.Database.Adapter)

	returnObject, err := testModel.Create(object)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	_, isTestObject := returnObject.(models.TestObject)
	if !isTestObject {
		srverr.Handle500(rw, fmt.Errorf("the return object is not a TestObject"))
		return
	}

	bytes, err := json.Marshal(returnObject)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Write(bytes)
}

func (c TestController) TestReadMethod(rw http.ResponseWriter, req *http.Request) {
	testModel := models.NewTestModel(&c.Database.Adapter)
	data, err := testModel.All()
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	testObjects, ok := data.([]models.TestObject)
	if !ok {
		srverr.Handle500(rw, fmt.Errorf("invalid data returned from the model"))
		return
	}

	bytes, err := json.Marshal(testObjects)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Write(bytes)
}

func (c TestController) TestUpdateMethod(rw http.ResponseWriter, req *http.Request) {
	var object models.TestObject
	err := util.DecodeRequestBody(req, &object)
	if err != nil {
		srverr.Handle400(rw, err)
		return
	}

	testModel := models.NewTestModel(&c.Database.Adapter)

	err = testModel.Update(object.Id, object)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Write([]byte("success"))
}

func (c TestController) TestDestroyMethod(rw http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")

	testModel := models.NewTestModel(&c.Database.Adapter)

	err := testModel.Delete(id)
	if err != nil {
		srverr.Handle500(rw, err)
		return
	}

	rw.Write([]byte("success"))
}

func (c TestController) CRUDTest(rw http.ResponseWriter, req *http.Request) {
	render.RenderView(rw, "test/crud_test.html", nil)
}
