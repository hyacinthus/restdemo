package main

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/go-redis/cache"
	"github.com/labstack/echo"
	"github.com/vmihailenco/msgpack"
)

// 初始化缓存
func initCache() {
	cc = &cache.Codec{
		Redis: rdb,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

// setcc 写缓存
func setcc(key string, object interface{}, exp time.Duration) {
	cc.Set(&cache.Item{
		Key:        key,
		Object:     object,
		Expiration: exp,
	})
}

// getcc 读缓存
func getcc(key string, pointer interface{}) error {
	return cc.Get(key, pointer)
}

// delcc 清缓存
func delcc(key string) {
	cc.Delete(key)
}

// cleancc 清除一类缓存
func cleancc(cate string) {
	if cate == "" {
		logrus.Error("someone try to clean all cache keys")
		return
	}
	i := 0
	for _, key := range rdb.Keys(cate + "*").Val() {
		delcc(key)
		i++
	}
	logrus.Infof("delete %d %s cache", i, cate)
}

func deleteCache(c echo.Context) error {
	cate := c.Param("cate")
	switch cate {
	case "token":
		cleancc("token")
	case "all":
		cleancc("token")
	default:
		return newHTTPError(400, "InvalidID", "请在URL中提供合法的缓存类型")
	}
	return c.NoContent(http.StatusNoContent)
}
