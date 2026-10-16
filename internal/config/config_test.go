package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	cfg := Load()
	if cfg.APIPort != "8080" {
		t.Errorf("APIPort = %q, esperado 8080", cfg.APIPort)
	}
	if cfg.DBPath != "/data/bgpadvert.db" {
		t.Errorf("DBPath = %q", cfg.DBPath)
	}
	if cfg.MaxPrefixes != "500" {
		t.Errorf("MaxPrefixes = %q", cfg.MaxPrefixes)
	}
	if cfg.APIKey != "" {
		t.Errorf("APIKey deveria virar vazio sem env, obtido %q", cfg.APIKey)
	}
}

func TestLoadOverride(t *testing.T) {
	t.Setenv("BGP_API_PORT", "9090")
	t.Setenv("BGP_DB_PATH", "/tmp/x.db")
	t.Setenv("BGP_API_KEY", "k1")
	cfg := Load()
	if cfg.APIPort != "9090" || cfg.DBPath != "/tmp/x.db" || cfg.APIKey != "k1" {
		t.Errorf("override falhou: %+v", cfg)
	}
}
