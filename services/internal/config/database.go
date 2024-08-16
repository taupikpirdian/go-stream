package config

import (
	"log"
	"os"
	"strings"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/lib/pq"
)

func InitDatabase() (*gorm.DB, error) {
	cfg := config.DatabasePGSQL()
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	configs := []string{
		"host=" + cfg.Host,
		"user=" + cfg.User,
		"password=" + cfg.Password,
		"dbname=" + cfg.Database,
		"port=" + cfg.Port,
		"TimeZone=" + cfg.TimeZone,
		"sslmode=disable",
		"search_path=" + cfg.Schema,
	}
	dsn := strings.Join(configs, " ")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxConnectionIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxConnectionOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}
