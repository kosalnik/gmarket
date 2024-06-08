package config

import (
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

//go:embed default/*.yaml
var configDir embed.FS

const (
	cfgDefaultPath = `default/default.yaml`
)

type Config struct {
	Server        Server
	Database      Database
	AccrualSystem AccrualSystem
	Logger        Logger
	JWT           JWT
}

type Server struct {
	Address string
}

type JWT struct {
	PublicKey  string
	PrivateKey string
}

type Database struct {
	URI string
}

type AccrualSystem struct {
	Address           string
	Concurrence       int
	RateLimit         int
	RateLimitDuration time.Duration
}

type Logger struct {
	Level string
}

func NewConfig() *Config {
	v := viper.NewWithOptions(viper.EnvKeyReplacer(&EnvKeyReplacer{}))
	//v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	defaultCfg, err := configDir.Open(cfgDefaultPath)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := v.ReadConfig(defaultCfg); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &cfg
}

type EnvKeyReplacer struct{}

var _ viper.StringReplacer = &EnvKeyReplacer{}

func (e *EnvKeyReplacer) Replace(key string) string {
	switch key {
	case "ACCRUALSYSTEM.ADDRESS":
		return "ACCRUAL_SYSTEM_ADDRESS"
	case "SERVER.ADDRESS":
		return "RUN_ADDRESS"
	}
	return strings.Replace(key, ".", "_", -1)
}
