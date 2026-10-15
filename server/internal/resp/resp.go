// Package resp implements the response envelope expected by the Arco Pro
// axios interceptor: code 20000 means success, anything else is an error.
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
	CodeBadRequest   = 40000
)

type Body struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
}

type Page[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Code: CodeOK, Status: "ok", Msg: "请求成功", Data: data})
}

func Fail(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(http.StatusOK, Body{Code: code, Status: "fail", Msg: msg})
}
