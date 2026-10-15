// Package resp writes responses in the envelope the Arco Pro frontend expects.
package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeOK           = 20000
	CodeError        = 50000
	CodeInvalidToken = 50008
	CodeTokenExpired = 50014
	CodeForbidden    = 40300
)

type Body struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Code   int    `json:"code"`
	Data   any    `json:"data"`
}

type Page struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Status: "ok", Msg: "请求成功", Code: CodeOK, Data: data})
}

// Fail keeps HTTP 200 so the frontend interceptor can surface msg to the user.
func Fail(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(http.StatusOK, Body{Status: "fail", Msg: msg, Code: code})
}
