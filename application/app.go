package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rexcfnghk/pricing-store/config"
)

type App struct {
	router    http.Handler
	rdb       *redis.Client
	tokenAuth *jwtauth.JWTAuth
}

func New(appConfig *config.AppConfig) *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{
			Addr:     appConfig.Datastore.Host,
			Username: appConfig.Datastore.Username,
			Password: appConfig.Datastore.Password,
		}),
		tokenAuth: jwtauth.New("HS256", []byte(appConfig.Jwt.Secret), nil),
	}
	app.loadRoutes()

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

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting server")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
