package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/cache"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

// LoginRequest 登录提供内容
type LoginRequest struct {
	// 用户名
	Username string `json:"username"`
	// 密码
	Password string `json:"password"`
}

// Token 以上第四步返回给客户端的token对象
type Token struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    int       `json:"-"`
}

// API状态 成功204 失败500
func getStatus(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// login 登录函数，demo中为了简洁就只有user password可以通过
// 实际应用中这个函数会相当复杂，要用正则判断输入的用户名是什么类型，然后调用相关函数去找用户。
// 还要兼容第三方登录，所以请求结构体也会更加复杂。
// @Tags 用户
// @Summary 登录
// @Description 用户登录
// @Accept  json
// @Produce  json
// @Param data body main.LoginRequest true "登录凭证"
// @Success 201 {object} main.Token
// @Failure 400 {object} main.httpError
// @Failure 401 {object} main.httpError
// @Failure 500 {object} main.httpError
// @Router /login [post]
func login(c echo.Context) error {
	// 判断何种方式登录，小程序为提供code
	var req = new(LoginRequest) // 输入请求
	if err := c.Bind(req); err != nil {
		return err
	}
	var t *Token
	if req.Username == "username" && req.Password == "password" {
		// 发行token
		t = &Token{
			Token:     newUUID(),
			ExpiresAt: time.Now().Add(time.Hour * 96),
			// 这个userid应该是检索出来的，这里为demo写死。
			UserID: 1,
		}
		setcc("token:"+t.Token, t, time.Hour*96)
	} else {
		return ErrAuthFailed
	}
	return c.JSON(http.StatusOK, t)
}

// skipper 这些不需要token
func skipper(c echo.Context) bool {
	method := c.Request().Method
	path := c.Path()
	// 先处理非GET方法，除了登录，现实中还可能有一些 webhooks
	switch path {
	case
		// 登录
		"/login":
		return true
	}
	// 从这里开始必须是GET方法
	if method != "GET" {
		return false
	}
	if path == "" {
		return true
	}
	resource := strings.Split(path, "/")[1]
	switch resource {
	case
		// 公开信息，把需要公开的资源每个一行写这里
		"swagger",
		"public":
		return true
	}
	return false
}

// Validator 校验token是否合法，顺便根据token在 context中赋值 user id
func validator(token string, c echo.Context) (bool, error) {
	// 调试后门
	logrus.Debug("token:", token)
	if config.APP.Debug && token == "debug" {
		c.Set("user_id", 1)
		return true, nil
	}
	// 寻找token
	var t = new(Token)
	err := getcc("token:"+token, t)
	if err == cache.ErrCacheMiss {
		return false, nil
	} else if err != nil {
		return false, err
	}
	// 设置用户
	c.Set("user_id", t.UserID)

	return true, nil
}

// 这个函数还有一种设计风格，就是只是返回userid，
// 以支持可选登录，在业务中判断userid如果是0就没有登录
func parseUser(c echo.Context) (userID int, err error) {
	userID, ok := c.Get("user_id").(int)
	if !ok || userID == 0 {
		return 0, ErrUnauthorized
	}
	return userID, nil
}
