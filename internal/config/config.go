package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

var validate = validator.New()

type SMTP struct {
	Host     string `yaml:"host" validate:"required,hostname|ip"`
	Port     int    `yaml:"port" validate:"required,min=1,max=65535"`
	Username string `yaml:"username" validate:"required,email"`
	Password string `yaml:"password" validate:"required"`
	From     string `yaml:"from" validate:"required,email"`
}

type Templates struct {
	TestEmailTemplate string `yaml:"test_email_template" validate:"required"`
}

type Data struct {
	CSVFilePath string `yaml:"csv_file_path" validate:"required"`
}

type Retry struct {
	MaxRetries      int `yaml:"max_retries" validate:"min=0,max=10"`
	InitialDelay    int `yaml:"initial_delay_seconds" validate:"min=1"` // in seconds
	BackoffMultiplier float64 `yaml:"backoff_multiplier" validate:"min=1"`
	MaxDelay        int `yaml:"max_delay_seconds" validate:"min=1"` // in seconds
}

type Config struct {
	SMTP      SMTP      `yaml:"smtp" validate:"required"`
	Templates Templates `yaml:"templates" validate:"required"`
	Data      Data      `yaml:"data" validate:"required"`
	Retry     Retry     `yaml:"retry" validate:"required"`
}

func LoadConfig() *Config {
	var cfg Config

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	configPath := filepath.Join("config", fmt.Sprintf("%s.yaml", env))

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("Config file does not exist", "path", configPath)
		os.Exit(1)
	}

	fileInfo, err := os.Stat(configPath)
	if err != nil {
		slog.Error("Failed to get config file info", "path", configPath, "error", err)
		os.Exit(1)
	}

	if fileInfo.Size() == 0 {
		slog.Error("Config file is empty", "path", configPath)
		os.Exit(1)
	}

	// Read config file
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("Failed to read config file", "path", configPath, "error", err)
		os.Exit(1)
	}

	// Validate config
	if err := validate.Struct(cfg); err != nil {
		slog.Error("Configuration validation failed", "error", err)
		os.Exit(1)
	}

	return &cfg
}