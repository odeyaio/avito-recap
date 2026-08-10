package app

import (
	"fmt"
	"log/slog"
	"time"

	"avito-recap/internal/engine"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `env:"APP_ENV" env-default:"development"`
	HTTP        HTTPConfig
	DatabaseURL string     `env:"DATABASE_URL" env-required:"true"`
	LogLevel    slog.Level `env:"LOG_LEVEL" env-default:"info"`
	Engine      engine.Config
}

type HTTPConfig struct {
	Address           string        `env:"HTTP_ADDRESS" env-default:":8080"`
	ShutdownTimeout   time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"10s"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" env-default:"10s"`
	ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"30s"`
	WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"30s"`
	IdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"120s"`
}

func LoadConfig() (Config, error) {
	const op = "app.LoadConfig"

	var config Config
	if err := cleanenv.ReadEnv(&config); err != nil {
		return Config{}, fmt.Errorf("%s: %w", op, err)
	}

	return config, nil
}
