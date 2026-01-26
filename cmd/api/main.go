package main

import (
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	"github.com/vukieuhaihoa/bookmark-management/internal/model"

	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
	"github.com/vukieuhaihoa/bookmark-management/pkg/jwtutils"
	"github.com/vukieuhaihoa/bookmark-management/pkg/logger"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-management/pkg/sqldb"
)

// @title Bookmark Management API
// @version 1.2
// @description This is the API documentation for the Bookmark Management service.
// @host localhost:8080
// @BasePath /
// @schemes http
// @SecurityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
//
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email vukieuhaihoa@gmail.com
func main() {
	logger.SetLogLevel()

	cfg, err := api.NewConfig()
	common.HandlerError(err)

	redisClient, err := redisPkg.NewClient("")
	common.HandlerError(err)

	dbClient, err := sqldb.NewClient("")
	common.HandlerError(err)

	dbClient.AutoMigrate(&model.User{})

	jwtGenerator, err := jwtutils.NewJWTGenerator("./private_key.pem")
	common.HandlerError(err)

	jwtValidator, err := jwtutils.NewJWTValidator("./public_key.pem")
	common.HandlerError(err)

	app := api.New(cfg, redisClient, dbClient, jwtGenerator, jwtValidator)
	app.Start()
}
