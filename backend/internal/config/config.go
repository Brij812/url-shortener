package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

type Config struct {
	Server struct {
		Port string `koanf:"port"`
	} `koanf:"server"`

	Database struct {
		Host     string `koanf:"host"`
		Port     string `koanf:"port"`
		User     string `koanf:"user"`
		Password string `koanf:"password"`
		Name     string `koanf:"name"`
		SSLMode  string `koanf:"sslmode"`
	} `koanf:"database"`

	Redis struct {
		Host     string `koanf:"host"`
		Port     string `koanf:"port"`
		Password string `koanf:"password"`
		DB       int    `koanf:"db"`
		PoolSize int    `koanf:"pool_size"`
	} `koanf:"redis"`

	JWT struct {
		Secret                   string `koanf:"secret"`
		Issuer                   string `koanf:"issuer"`
		AccessTokenExpiryMinutes int    `koanf:"access_token_expiry_minutes"`
		RefreshTokenExpiryHours  int    `koanf:"refresh_token_expiry_hours"`
	} `koanf:"jwt"`
}

func LoadConfig() *Config {
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		log.Println("⚠️ No config.yaml found, skipping file load")
	}

	k.Load(env.Provider("", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(s), "_", ".")
	}), nil)

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return &cfg
}
