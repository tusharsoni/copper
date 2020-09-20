package csql

import (
	"fmt"

	"gorm.io/gorm/logger"

	"github.com/tusharsoni/copper/clogger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDBParams struct {
	Config Config
	Logger clogger.Logger
}

func NewGormDB(p GormDBParams) (*gorm.DB, error) {
	conn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		p.Config.Host,
		p.Config.Port,
		p.Config.User,
		p.Config.Name,
	)

	p.Logger.WithTags(map[string]interface{}{
		"connection": conn,
	}).Info("Connecting to database...")

	if p.Config.Password != "" {
		conn = fmt.Sprintf("%s password=%s", conn, p.Config.Password)
	}

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		p.Logger.Error("Failed to connect to database..", err)
		return nil, err
	}

	return db, nil
}
