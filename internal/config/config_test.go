package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("recover %s", err)
		}
	}()

	t.Run("Config: Default", func(t *testing.T) {
		cfg := NewConfig()
		if cfg == nil {
			t.Errorf(`Want Config instance, got nil`)
		}
	})

	t.Run("Config: ENV rewrite", func(t *testing.T) {
		withEnv := func(envs map[string]string, cb func()) {
			old := make(map[string]string, len(envs))
			for k := range envs {
				old[k] = os.Getenv(k)
			}
			defer func() {
				for k, v := range old {
					require.NoError(t, os.Setenv(k, v))
				}
			}()
			for k, v := range envs {
				require.NoError(t, os.Setenv(k, v))
			}
			cb()
		}
		withEnv(map[string]string{
			"DATABASE_URI":           "asd",
			"RUN_ADDRESS":            "qwe",
			"ACCRUAL_SYSTEM_ADDRESS": "zxc",
		}, func() {
			cfg := NewConfig()
			if cfg == nil {
				t.Errorf(`Want Config instance, got nil`)
			} else { // Статанализатор ругается, если следующие ассерты вне if вставить
				assert.Equal(t, "asd", cfg.Database.URI)
				assert.Equal(t, "qwe", cfg.Server.Address)
				assert.Equal(t, "zxc", cfg.AccrualSystem.Address)
			}
		})
	})
}
