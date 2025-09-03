package config

import (
	"fmt"

	"github.com/artyomkorchagin/first-task/pkg/helpers"
)

func GetDSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		helpers.GetEnv("DB_HOST", ""),
		helpers.GetEnv("DB_PORT", ""),
		helpers.GetEnv("DB_NAME", ""),
		helpers.GetEnv("DB_USER", ""),
		helpers.GetEnv("DB_PASSWORD", ""),
		helpers.GetEnv("DB_SSLMODE", ""))
}

func GetDriver() string {
	return helpers.GetEnv("DB_DRIVER", "pgx")
}
