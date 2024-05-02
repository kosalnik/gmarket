package config

type Config struct {
	Server        Server
	Database      Database
	AccrualSystem AccrualSystem
	Logger        Logger
}

type Server struct {
	Address string
}

type Database struct {
	URI string
}

type AccrualSystem struct {
	Address string
}

type Logger struct {
	Level string
}

func NewConfig() *Config {
	return &Config{
		Server:        Server{Address: ":8080"},
		Database:      Database{URI: ""},
		AccrualSystem: AccrualSystem{Address: ""},
		Logger:        Logger{Level: "info"},
	}
}
