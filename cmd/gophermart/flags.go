package main

import (
	"flag"

	"github.com/kosalnik/gmarket/internal/config"
)

func parseFlags(c *config.Config) {
	flag.StringVar(&c.Server.Address, "a", ":8080", "server endpoint (ip:port)")
	flag.StringVar(&c.AccrualSystem.Address, "r", "", "Accrual system address")
	flag.StringVar(&c.Database.URI, "d", "", "Database DSN")
	flag.Parse()
}
