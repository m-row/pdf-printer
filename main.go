package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

var (
	e = echo.New()
	c = Controller{s}
	s = GetSettings()

	CommitCount    = "0"
	CommitDescribe = "dev"
	Version        = fmt.Sprintf(
		"%s.%s.%s",
		s.AppVer,
		CommitCount,
		CommitDescribe,
	)
)

func main() {
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	e.GET("/meta", c.Meta).Name = "meta"
	e.POST("/print", c.Print).Name = "print"

	if err := serve(e); err != nil {
		slog.Error(fmt.Sprintf("Serve error: %s", err.Error()))
	}
}

func serve(e *echo.Echo) error {
	srv := &http.Server{
		Addr:        fmt.Sprintf(":%s", s.Port),
		Handler:     e,
		IdleTimeout: time.Minute,
		// ErrorLog:     log.New(app.Logger, "", 0),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	slog.Info(fmt.Sprintf("starting server at http://localhost%s", srv.Addr))

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		slog.Info(fmt.Sprintf("shutting down server, signal: %s", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("stopped server, Addr: http://localhost%s", srv.Addr))

	return nil
}
