package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dailoi280702/vrs-ranking-service/client/generalservice"
	_ "github.com/dailoi280702/vrs-ranking-service/client/mysql"
	_ "github.com/dailoi280702/vrs-ranking-service/client/redis"
	"github.com/dailoi280702/vrs-ranking-service/config"
	_ "github.com/dailoi280702/vrs-ranking-service/docs"
	"github.com/dailoi280702/vrs-ranking-service/handler/http"
	"github.com/dailoi280702/vrs-ranking-service/log"
)

// @title			Swagger Video ranking service API
// @version		1.0
// @description	Server for using Swagger with Echo.
// @host			localhost:9000
// @BasePath		/api/v1
func main() {
	cfg := config.GetConfig()

	var (
		logger = log.Logger()
		errs   = make(chan error)
		h      = http.NewHTTPHandler()
	)

	defer generalservice.Close()

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
