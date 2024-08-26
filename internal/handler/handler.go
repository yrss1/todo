package handler

import (
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/todo/internal/config"
	"github.com/yrss1/todo/internal/handler/http"
	"github.com/yrss1/todo/internal/service/account"
	"github.com/yrss1/todo/internal/service/todo"
	"github.com/yrss1/todo/pkg/server/response"
	"github.com/yrss1/todo/pkg/server/router"
)

type Dependencies struct {
	Configs config.Configs

	AccountService *account.Service
	TodoService    *todo.Service
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

		//docs.SwaggerInfo.BasePath = h.dependencies.Configs.APP.Path
		//h.HTTP.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		userHandler := http.NewUserHandler(h.dependencies.AccountService)
		taskHandler := http.NewTaskHandler(h.dependencies.TodoService)

		api := h.HTTP.Group(h.dependencies.Configs.APP.Path)
		{
			userHandler.Routes(api)
			taskHandler.Routes(api)
		}
		return
	}
}
