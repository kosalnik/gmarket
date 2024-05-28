package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("recover %s", err)
		}
	}()
	cfg := NewConfig()
	if cfg == nil {
		t.Errorf(`Want Config instance, got nil`)
	}
}
