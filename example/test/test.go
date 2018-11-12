package test

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"mime/multipart"

	"github.com/wusphinx/gin-swagger/example/from_request"
	"github.com/wusphinx/gin-swagger/example/globals"
	"github.com/wusphinx/gin-swagger/example/test2"
)

// ErrorMap
type ErrorMap map[string]map[string]*int

// SomeTest
type SomeTest struct {
	test2.Common
	State    test2.State `json:"state" validate:"@string{TWO}"`
	ErrorMap ErrorMap    `json:"errorMap"`
}

// ReqBody
type ReqBody struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
}

// SomeReq
type SomeReq struct {
	// Body
	test2.Pager
	StartTime test2.Date           `in:"query" json:"startTime"`
	State     test2.State          `in:"query" json:"state" validate:"@string{TWO}"`
	File      multipart.FileHeader `in:"formData" json:"file"`
}

// @httpError(40000200,HTTP_ERROR_UNKNOWN,"未定义","",false);
func someDoReq() {
}

func Test(c *gin.Context) {
	req := SomeReq{}

	fmt.Println(req)

	var res = SomeTest{
		State: test2.STATE__ONE,
	}

	someDoReq()
	globals.WriteErr(c)

	// 正常返回
	c.JSON(http.StatusOK, res)
}

type AuthReq struct {
	Authorization string `json:"authorization" in:"header"`
}

func (req AuthReq) Handle(c *gin.Context) {
	if req.Authorization == "" {
		c.JSON(globals.HTTP_ERROR_UNKNOWN.Status(), globals.HTTP_ERROR_UNKNOWN.ToError())
	}
}

func AuthMiddleware(c *gin.Context) {
	var req = AuthReq{}
	if req.Authorization == "" {
		c.JSON(globals.HTTP_ERROR_UNKNOWN.Status(), globals.HTTP_ERROR_UNKNOWN.ToError())
	}
}

func Auth() gin.HandlerFunc {
	return from_request.FromRequest(AuthReq{})
}
