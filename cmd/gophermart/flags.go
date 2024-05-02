package main

import (
	"flag"
	"os"

	"github.com/kosalnik/gmarket/internal/config"
)

func parseFlags(c *config.Config) {
	flag.StringVar(&c.Server.Address, "a", ":8080", "server endpoint (ip:port)")
	flag.StringVar(&c.AccrualSystem.Address, "r", "", "Accrual system address")
	flag.StringVar(&c.Database.URI, "d", "", "Database DSN")
	flag.Parse()

	if v := os.Getenv("RUN_ADDRESS"); v != "" {
		c.Server.Address = v
	}
	if v := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); v != "" {
		c.AccrualSystem.Address = v
	}
	if v := os.Getenv("DATABASE_URI"); v != "" {
		c.Database.URI = v
	}
}
