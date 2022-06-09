// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/kordape/tweety/config"
	v1 "github.com/kordape/tweety/internal/controller/http/v1"
	"github.com/kordape/tweety/internal/tweets"
	"github.com/kordape/tweety/internal/tweets/webapi"
	"github.com/kordape/tweety/pkg/httpserver"
	"github.com/kordape/tweety/pkg/logger"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	// Use case
	tweetsClassifier := tweets.NewClassfier(
		webapi.New(
			cfg.TwitterAccessKey,
			cfg.TwitterSecretKey,
			cfg.TwitterBearerToken,
		),
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, log, tweetsClassifier)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		// Shutdown
		err = httpServer.Shutdown()
		if err != nil {
			log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	}
}
