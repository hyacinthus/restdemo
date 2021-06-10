package main

import (
	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// 定义全区变量 为了保证执行顺序 初始化均在main中执行
var (
	// gorm mysql db connection
	db *gorm.DB
	// redis client
	rdb *redis.Client
	// global cache
	cc *cache.Codec
)

// @title RESTful API DEMO by Golang & Echo
// @version 1.0
// @description This is a demo server.

// @contact.name Muninn
// @contact.email hyacinthus@gmail.com

// @license.name MIT
// @license.url https://github.com/hyacinthus/restdemo/blob/master/LICENSE

// @host demo.crandom.com
// @BasePath /
func main() {
	// init echo
	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:   skipper,   // 跳过验证条件 在 auth.go 定义
		Validator: validator, // 处理验证结果 在 auth.go 定义
	}))
	e.Use(ParsePagination) // 分页参数解析，在 pagination.go 定义

	// Echo debug setting
	if config.APP.Debug {
		e.Debug = true
	}

	// init mysql and redis
	initDB()
	defer db.Close()
	initRedis()
	defer rdb.Close()

	// init global cache
	initCache()

	// async create tables
	go createTables()

	// status
	e.GET("/status", getStatus)

	// auth
	e.POST("/login", login)

	// note Routes
	e.GET("/notes", getNotes)
	e.POST("/notes", createNote)
	e.GET("/notes/:id", getNote)
	e.PUT("/notes/:id", updateNote)
	e.DELETE("/notes/:id", deleteNote)
	e.GET("/public/notes", getPublicNotes)
	e.GET("/public/notes/:id", getPublicNote)

	// Start echo server
	e.Logger.Fatal(e.Start(config.APP.Host + ":" + config.APP.Port))
}
