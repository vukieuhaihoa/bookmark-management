package main

import (
	"github.com/vukieuhaihoa/bookmark-management/internal/api"

	"github.com/vukieuhaihoa/bookmark-management/pkg/logger"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
)

// @title Bookmark Management API
// @version 1.2
// @description This is the API documentation for the Bookmark Management service.
// @host localhost:8080
// @BasePath /
// @schemes http

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email author@gmail.com
func main() {
	logger.SetLogLevel()

	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}

	app := api.New(cfg, redisClient)
	app.Start()
}
