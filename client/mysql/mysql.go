package mysql

import (
	"fmt"

	"github.com/dailoi280702/vrs-ranking-service/config"
	"github.com/dailoi280702/vrs-ranking-service/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var (
		err    error
		cfg    = config.GetConfig()
		logger = log.Logger()
	)

	connectionString := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s",
		cfg.MySQL.User,
		cfg.MySQL.Pass,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DBName,
	)

	db, err = gorm.Open(mysql.New(mysql.Config{DSN: connectionString}), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to MySQL", "error", err, "config", cfg.MySQL)
	}

	logger.Info("Connected to mysql")
}

func GetClient() *gorm.DB {
	return &gorm.DB{
		Config: db.Config,
	}
}
