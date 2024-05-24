package app

import (
	"EffectiveMobile/config"
	v1 "EffectiveMobile/internal/controller/http/v1"
	"EffectiveMobile/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	l *zap.Logger
	c config.Config
}

func NewApp(log *zap.Logger, config config.Config) (*App, error) {
	return &App{
		l: log,
		c: config,
	}, nil
}

func (a *App) Run() {

	a.l.Debug("Debug mode is on")

	if !a.c.DEBUG {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	storage := postgres.NewStorage(&a.c)
	service := postgres.NewService(storage.Db)
	a.l.Debug("Postgres service created")

	v1.NewRouter(r, a.l, *service)

	// Graceful shutdown mechanism
	addr := fmt.Sprintf("%s:%s", a.c.HTTP.Host, a.c.HTTP.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		a.l.Info(fmt.Sprintf("Starting on %s...", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.l.Fatal("Could not listen", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.l.Warn("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.l.Fatal("Server forced to shutdown:", zap.Error(err))
	}
	// End of graceful shutdown mechanism

}
