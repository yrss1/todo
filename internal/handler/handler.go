package handler

import (
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yrss1/todo/docs"
	"github.com/yrss1/todo/internal/config"
	"github.com/yrss1/todo/internal/handler/http"
	"github.com/yrss1/todo/internal/service/account"
	"github.com/yrss1/todo/internal/service/auth"
	"github.com/yrss1/todo/internal/service/todo"
	"github.com/yrss1/todo/pkg/server/response"
	"github.com/yrss1/todo/pkg/server/router"
)

type Dependencies struct {
	Configs config.Configs

	AccountService *account.Service
	TodoService    *todo.Service
	AuthService    *auth.Service
}
type Handler struct {
	dependencies Dependencies
	HTTP         *gin.Engine
}
type Configuration func(h *Handler) error

func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	h = &Handler{
		dependencies: d,
	}

	for _, cfg := range configs {
		if err = cfg(h); err != nil {
			return
		}
	}

	return
}

// @title Todo API
// @version 1.0
// @description This is a sample server for a todo application.
// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func WithHTTPHandler() Configuration {
	return func(h *Handler) (err error) {
		h.HTTP = router.New()
		h.HTTP.Use(timeout.New(
			timeout.WithTimeout(h.dependencies.Configs.APP.Timeout),
			timeout.WithHandler(func(ctx *gin.Context) {
				ctx.Next()
			}),
			timeout.WithResponse(func(ctx *gin.Context) {
				response.StatusRequestTimeout(ctx)
			}),
		))

		docs.SwaggerInfo.BasePath = h.dependencies.Configs.APP.Path
		h.HTTP.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		authHandler := http.NewAuthHandler(h.dependencies.AuthService, h.dependencies.Configs.APP.JWT)
		healthHandler := http.NewHealthHandler()

		authAPI := h.HTTP.Group(h.dependencies.Configs.APP.Path)
		{
			authHandler.Routes(authAPI)
			healthHandler.Routes(authAPI)
		}

		userHandler := http.NewUserHandler(h.dependencies.AccountService)
		taskHandler := http.NewTaskHandler(h.dependencies.TodoService)

		api := h.HTTP.Group(h.dependencies.Configs.APP.Path)
		{
			api.Use(authHandler.AuthMiddleware())

			userHandler.Routes(api)
			taskHandler.Routes(api)
		}
		return
	}
}
