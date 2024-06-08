package config

import (
	"embed"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/pflag"
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
	parseFlags(v)
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

func parseFlags(v *viper.Viper) {
	flag.String("a", ":8080", "server endpoint (ip:port)")
	flag.String("r", "", "Accrual system address")
	flag.String("d", "", "Database DSN")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := v.BindPFlag("Server.Address", pflag.Lookup("a")); err != nil {
		panic(err)
	}
	if err := v.BindPFlag("AccrualSystem.Address", pflag.Lookup("r")); err != nil {
		panic(err)
	}
	if err := v.BindPFlag("Database.URI", pflag.Lookup("d")); err != nil {
		panic(err)
	}
}

type EnvKeyReplacer struct{}

var _ viper.StringReplacer = &EnvKeyReplacer{}

func (e *EnvKeyReplacer) Replace(key string) string {
	switch key {
	case "ACCRUALSYSTEM.ADDRESS":
		return "ACCRUAL_SYSTEM_ADDRESS"
	case "SERVER.ADDRESS":
		return "RUN_ADDRESS"
	case "DATABASE.URI":
		return "DATABASE_URI"
	}
	return strings.Replace(key, ".", "_", -1)
}
