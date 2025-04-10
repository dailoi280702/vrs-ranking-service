package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/dailoi280702/vrs-ranking-service/client/mysql"
	_ "github.com/dailoi280702/vrs-ranking-service/client/redis"
	"github.com/dailoi280702/vrs-ranking-service/config"
	"github.com/dailoi280702/vrs-ranking-service/handler/http"
	"github.com/dailoi280702/vrs-ranking-service/log"
)

func main() {
	cfg := config.GetConfig()

	var (
		logger = log.Logger()
		errs   = make(chan error)
		h      = http.NewHTTPHandler()
	)

	go func() {
		logger.Info("Server is running", "port", cfg.Port)
		errs <- h.Start(":" + cfg.Port)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Info("The server is stopping ...")

	_ = h.Shutdown(context.Background())

	close(c)

	logger.Info("Exited", "error", <-errs)
}
