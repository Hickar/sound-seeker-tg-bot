package postgres

import (
	"fmt"
	"strings"

	"github.com/Hickar/sound-seeker-bot/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(conf config.DatabaseConfig) (*gorm.DB, error) {
	hostParts := strings.Split(conf.Host, ":")
	if len(hostParts) != 2 {
		return nil, fmt.Errorf("malformed postgres host: %s", conf.Host)
	}
	host, port := hostParts[0], hostParts[1]

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s", host, conf.User, conf.Password, port, conf.Name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}

	return db, nil
}