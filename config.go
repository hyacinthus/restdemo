package main

import (
	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var config = struct {
	APP struct {
		Debug    bool   `default:"false"`
		Host     string `default:"0.0.0.0"`
		Port     string `default:"1324"`
		PageSize int    `default:"10"`
		BaseURL  string `default:"https://api.example.com/"`
		FileURL  string `default:"https://static.example.com/"`
	}

	DB struct {
		Host     string `default:"mysql"`
		Port     string `default:"3306"`
		User     string `default:"root"`
		Password string `default:"root"`
		Name     string `default:"art"`
	}

	Redis struct {
		Host     string `default:"redis"`
		Port     string `default:"6379"`
		Password string
		DB       int `default:"0"`
	}
}{}

func init() {
	godotenv.Load()
	configor.Load(&config)
	if config.APP.Debug {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "06-01-02 15:04:05.00",
		})
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
