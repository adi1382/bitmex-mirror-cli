package config

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	config := LoadConfig("../config.txt")
	if config.Host == "" {
		t.Error("No config")
	}
	if config.Key == "" {
		t.Error("No config")
	}
	if config.Secret == "" {
		t.Error("No config")
	}
}

func TestGetMasterConfig(t *testing.T) {
	config := LoadMasterConfig("../config.txt")
	if config.Master.Host == "" {
		t.Error("No config")
	}
}
