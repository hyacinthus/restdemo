package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// FileURL 图片链接
type FileURL string

// ToString 转换为string类型
func (f FileURL) ToString() string {
	var s = string(f)
	var url = s
	if !strings.HasPrefix(s, "http") {
		url = config.APP.FileURL + s
	}
	return url
}

// MarshalJSON 转换为json类型 加域名
func (f FileURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.ToString())
}

// UnmarshalJSON 不做处理
func (f *FileURL) UnmarshalJSON(data []byte) error {
	var tmp string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	tmp = strings.TrimPrefix(tmp, config.APP.FileURL)
	*f = FileURL(tmp)
	return nil
}

// Scan implements the Scanner interface.
func (f *FileURL) Scan(src interface{}) error {
	if src == nil {
		*f = ""
		return nil
	}
	tmp, ok := src.([]byte)
	if !ok {
		return errors.New("Read file url data from DB failed")
	}
	*f = FileURL(tmp)
	return nil
}

// Value implements the driver Valuer interface.
func (f FileURL) Value() (driver.Value, error) {
	return string(f), nil
}

func newUUID() string {
	id, err := uuid.NewV1()
	if err != nil {
		logrus.Errorf("Gen uuid Error:%s", err.Error())
		return "error"
	}
	return id.String()
}
