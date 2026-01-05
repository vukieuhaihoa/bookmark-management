// Package api provides the HTTP API server implementation for the bookmark management application.
// It handles routing, middleware, and server lifecycle management using the Gin web framework.
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/vukieuhaihoa/bookmark-management/docs"
	"github.com/vukieuhaihoa/bookmark-management/internal/handler"
	"github.com/vukieuhaihoa/bookmark-management/internal/repository"
	"github.com/vukieuhaihoa/bookmark-management/internal/service"
	"github.com/vukieuhaihoa/bookmark-management/pkg/stringutils"
)

// Engine defines the contract for the HTTP server engine.
// It provides methods to start and manage the API server lifecycle.
type Engine interface {
	// Start initializes and starts the HTTP server.
	// Returns an error if the server fails to start or encounters a runtime error.
	Start() error

	// ServeHTTP allows the engine to handle HTTP requests.
	// It satisfies the http.Handler interface.
	// Parameters:
	//   - w: The http.ResponseWriter to write the response
	//   - r: The incoming http.Request to be handled
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// api is the concrete implementation of the Engine interface.
// It encapsulates the Gin engine and provides HTTP server functionality.
type api struct {
	// app is the underlying Gin engine instance that handles HTTP routing and middleware
	app *gin.Engine

	// cfg holds the configuration settings for the API server
	cfg *Config

	// redisClient is the Redis client used for caching and session management
	redisClient *redis.Client

	// randomCodeGen is the code generator used for generating random codes
	randomCodeGen stringutils.CodeGenerator
}

// New creates a new instance of the API server engine.
// It initializes the Gin engine, configures routes, and prepares the server for handling requests.
//
// Parameters:
//   - cfg: The configuration settings for the API server
//   - redisClient: The Redis client used for caching and session management
//
// Returns:
//   - Engine: The initialized API server engine instance
func New(cfg *Config, redisClient *redis.Client) Engine {
	a := &api{
		app:           gin.New(),
		cfg:           cfg,
		redisClient:   redisClient,
		randomCodeGen: stringutils.NewCodeGenerator(),
	}

	a.registerRoutes()
	return a
}

// Start initializes and starts the HTTP server.
// It listens on the configured application port and begins handling incoming requests.
// defaults to port 8080 if not specified in the configuration.
// Returns:
//   - error: An error if the server fails to start, nil otherwise
func (a *api) Start() error {
	return a.app.Run(a.cfg.AppPort)
}

// ServeHTTP allows the api struct to satisfy the http.Handler interface.
// It delegates HTTP requests to the underlying Gin engine.
//
// Parameters:
//   - w: The http.ResponseWriter to write the response
//   - r: The incoming http.Request to be handled
func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

// registerRoutes configures all HTTP routes and their handlers for the API.
// This method initializes services and handlers, then registers them with the Gin engine.
// Currently registers the password generation endpoint.
func (a *api) registerRoutes() {
	healthCheckRepo := repository.NewHealthCheck(a.redisClient)

	passSvc := service.NewPassword(a.randomCodeGen)
	passHandler := handler.NewPassword(passSvc)

	healthCheckSvc := service.NewHealthCheck(a.cfg.ServiceName, a.cfg.InstanceID, healthCheckRepo)
	healthCheckHandler := handler.NewHealthCheck(healthCheckSvc)

	shortenURLRepo := repository.NewUrlStorage(a.redisClient)
	shortenURLSvc := service.NewShortenURL(shortenURLRepo, a.randomCodeGen)
	shortenURLHandler := handler.NewShortenURL(shortenURLSvc)

	// API version group
	v1 := a.app.Group("/v1")
	{
		// Register password generation endpoint
		v1.GET("/generate-password", passHandler.GeneratePassword)

		// Register URL shortening endpoint
		v1.POST("/links/shorten", shortenURLHandler.ShortenURL)
	}

	// Register health check endpoint
	a.app.GET("/health-check", healthCheckHandler.Check)

	// Register Swagger documentation endpoint
	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
