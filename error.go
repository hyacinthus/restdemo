package main

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

// 定义错误
var (
	ErrNotFound     = newHTTPError(404, "NotFound", "没有找到相应记录")
	ErrAuthFailed   = newHTTPError(401, "AuthFailed", "登录失败")
	ErrUnauthorized = newHTTPError(401, "Unauthorized", "本接口只有登录用户才能调用")
	ErrForbidden    = newHTTPError(403, "Forbidden", "权限不足")
)

// httpError 对外输出的错误格式
type httpError struct {
	code int
	// 错误代码，为英文字符串，前端可用此判断大的错误类型。
	Key string `json:"error"`
	// 错误消息，为详细错误描述，前端可选择性的展示此字段。
	Message string `json:"message"`
}

func newHTTPError(code int, key string, msg string) *httpError {
	return &httpError{
		code:    code,
		Key:     key,
		Message: msg,
	}
}

// Error makes it compatible with `error` interface.
func (e *httpError) Error() string {
	return e.Key + ": " + e.Message
}

// httpErrorHandler customize echo's HTTP error handler.
func httpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		key  = "ServerError"
		msg  string
	)
	// 二话不说先打日志
	c.Logger().Error(err.Error())

	if he, ok := err.(*httpError); ok {
		// 我们自定的错误
		code = he.code
		key = he.Key
		msg = he.Message
	} else if ee, ok := err.(*echo.HTTPError); ok {
		// echo 框架的错误
		code = ee.Code
		key = http.StatusText(code)
		msg = key
	} else if err == gorm.ErrRecordNotFound {
		// 我们将 gorm 的没有找到直接返回 404
		code = http.StatusNotFound
		key = "NotFound"
		msg = "没有找到相应记录"
	} else if config.APP.Debug {
		// 剩下的都是500 开了debug显示详细错误
		msg = err.Error()
	} else {
		// 500 不开debug 用标准错误描述 以防泄漏信息
		msg = http.StatusText(code)
	}

	// 判断 context 是否已经返回了
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err := c.NoContent(code)
			if err != nil {
				c.Logger().Error(err.Error())
			}
		} else {
			err := c.JSON(code, newHTTPError(code, key, msg))
			if err != nil {
				c.Logger().Error(err.Error())
			}
		}
	}
}
