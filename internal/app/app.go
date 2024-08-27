package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/yrss1/todo/internal/config"
	"github.com/yrss1/todo/internal/handler"
	"github.com/yrss1/todo/internal/repository"
	"github.com/yrss1/todo/internal/service/account"
	"github.com/yrss1/todo/internal/service/auth"
	"github.com/yrss1/todo/internal/service/todo"
	"github.com/yrss1/todo/pkg/log"
	"github.com/yrss1/todo/pkg/server"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	logger := log.LoggerFromContext(context.Background())

	configs, err := config.New()
	if err != nil {
		logger.Error("ERR_INIT_CONFIGS", zap.Error(err))
		return
	}

	repositories, err := repository.New(repository.WithPostgresStore(configs.POSTGRES.DSN))
	if err != nil {
		logger.Error("ERR_INIT_REPOSITORIES", zap.Error(err))
		return
	}

	authService, err := auth.New(
		auth.WithUserRepository(repositories.User),
	)

	accountService, err := account.New(
		account.WithUserRepository(repositories.User),
	)
	if err != nil {
		logger.Error("ERR_INIT_ACCOUNT_SERVICE", zap.Error(err))
		return
	}

	todoService, err := todo.New(
		todo.WithTaskRepository(repositories.Task))
	if err != nil {
		logger.Error("ERR_INIT_TODO_SERVICE", zap.Error(err))
		return
	}

	handlers, err := handler.New(
		handler.Dependencies{
			Configs:        configs,
			AuthService:    authService,
			AccountService: accountService,
			TodoService:    todoService,
		},
		handler.WithHTTPHandler())
	if err != nil {
		logger.Error("ERR_INIT_HANDLERS", zap.Error(err))
		return
	}

	servers, err := server.New(server.WithHTTPServer(handlers.HTTP, configs.APP.Port))
	if err != nil {
		logger.Error("ERR_INIT_SERVERS", zap.Error(err))
		return
	}
	if err = servers.Run(); err != nil {
		logger.Error("ERR_RUN_SERVERS", zap.Error(err))
		return
	}
	logger.Info("http server started on http://localhost:" + configs.APP.Port + "/swagger/index.html")

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	if err = servers.Stop(ctx); err != nil {
		panic(err)
	}

	fmt.Println("running cleanup tasks...")

	fmt.Println("server was successful shutdown.")
}
