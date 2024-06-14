package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer       ServerCfg      `yaml:"http_server"`
	Database         DatabaseConfig `yaml:"database"`
	JWT              JWTCfg         `yaml:"auth"`
	DefaultAdminPass string         `yaml:"default_admin_pass"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type ServerCfg struct {
	Addr        string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"120s"`
}

type JWTCfg struct {
	Secret string `yaml:"secret"`
}

type ConfigSMTP struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

func MustLoad() *Config {
	var cfg Config

	err := cleanenv.ReadConfig("../config/config.yaml", &cfg)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
