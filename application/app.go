package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/rexcfnghk/pricing-store/config"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New(appConfig *config.AppConfig) *App {
	app := &App{
		router: loadRoutes(),
		rdb: redis.NewClient(&redis.Options{
			Addr:     appConfig.Datastore.Host,
			Username: appConfig.Datastore.Username,
			Password: appConfig.Datastore.Password,
		}),
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	err = server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
