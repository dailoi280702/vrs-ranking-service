package generalservice

import (
	"github.com/dailoi280702/vrs-general-service/proto"
	"github.com/dailoi280702/vrs-ranking-service/config"
	"github.com/dailoi280702/vrs-ranking-service/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	conn   *grpc.ClientConn
	client proto.ServiceClient
)

func init() {
	var (
		err    error
		cfg    = config.GetConfig()
		logger = log.Logger()
	)

	conn, err = grpc.NewClient(
		cfg.GeneralServiceEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("Failed to connect general service", "error", err, "endpoint", cfg.GeneralServiceEndpoint)
	} else {
		logger.Info("Connected to general service", "host", cfg.GeneralServiceEndpoint)
	}

	client = proto.NewServiceClient(conn)
}

func Close() {
	if conn == nil {
		return
	}

	if err := conn.Close(); err != nil {
		log.Logger().Error("Failed to disconnect general service", "error", err)
	}
}

func GetClient() proto.ServiceClient {
	return client
}
