package config

import (
	"fmt"
	"log"
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

type Config struct {
	SMTP      SMTP      `yaml:"smtp" validate:"required"`
	Templates Templates `yaml:"templates" validate:"required"`
	Data      Data      `yaml:"data" validate:"required"`
}

func LoadConfig() *Config {
	var cfg Config

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	configPath := filepath.Join("config", fmt.Sprintf("%s.yaml", env))

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file %s: %v", configPath, err)
	}

	if err := validate.Struct(cfg); err != nil {
		log.Fatalf("configuration validation failed: %v", err)
	}

	return &cfg
}