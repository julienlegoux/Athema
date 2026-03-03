package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")

	yamlContent := `
server:
  host: "localhost"
  port: 8080
db:
  host: "localhost"
  port: 5432
  database: "athema"
  username: "athema"
  password: "secret"
  sslmode: "disable"
log:
  level: "info"
  format: "json"
memory:
  enabled: true
conversation:
  enabled: true
personality:
  enabled: true
emotional:
  enabled: true
lifecycle:
  enabled: true
initiation:
  enabled: true
`
	if err := os.WriteFile(cfgPath, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("Server.Host = %q, want %q", cfg.Server.Host, "localhost")
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 8080)
	}
	if cfg.DB.Database != "athema" {
		t.Errorf("DB.Database = %q, want %q", cfg.DB.Database, "athema")
	}
	if cfg.DB.Port != 5432 {
		t.Errorf("DB.Port = %d, want %d", cfg.DB.Port, 5432)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("Log.Level = %q, want %q", cfg.Log.Level, "info")
	}
	if !cfg.Memory.Enabled {
		t.Error("Memory.Enabled = false, want true")
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")

	yamlContent := `
server:
  host: "localhost"
  port: 8080
db:
  host: "localhost"
  port: 5432
  database: "athema"
  username: "athema"
  password: "secret"
  sslmode: "disable"
log:
  level: "info"
  format: "json"
`
	if err := os.WriteFile(cfgPath, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("ATHEMA_SERVER_PORT", "9090")
	t.Setenv("ATHEMA_DB_HOST", "db.prod.example.com")
	t.Setenv("ATHEMA_LOG_LEVEL", "debug")

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Server.Port = %d, want %d (env override)", cfg.Server.Port, 9090)
	}
	if cfg.DB.Host != "db.prod.example.com" {
		t.Errorf("DB.Host = %q, want %q (env override)", cfg.DB.Host, "db.prod.example.com")
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %q, want %q (env override)", cfg.Log.Level, "debug")
	}
	// Non-overridden values should remain from YAML
	if cfg.Server.Host != "localhost" {
		t.Errorf("Server.Host = %q, want %q (YAML default)", cfg.Server.Host, "localhost")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Error("Load() expected error for nonexistent file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(cfgPath, []byte("{{invalid yaml"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(cfgPath)
	if err == nil {
		t.Error("Load() expected error for invalid YAML, got nil")
	}
}

func TestDBConfig_DSN(t *testing.T) {
	cfg := DBConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "athema",
		Username: "user",
		Password: "pass",
		SSLMode:  "disable",
	}
	want := "postgres://user:pass@localhost:5432/athema?sslmode=disable"
	if got := cfg.DSN(); got != want {
		t.Errorf("DSN() = %q, want %q", got, want)
	}
}

func TestDBConfig_DSN_SpecialChars(t *testing.T) {
	cfg := DBConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "athema",
		Username: "user@org",
		Password: "p@ss:word",
		SSLMode:  "disable",
	}
	got := cfg.DSN()
	want := "postgres://user%40org:p%40ss%3Aword@localhost:5432/athema?sslmode=disable"
	if got != want {
		t.Errorf("DSN() = %q, want %q", got, want)
	}
}

func TestLoad_SubsystemEnvOverride(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")

	yamlContent := `
server:
  host: "localhost"
  port: 8080
db:
  host: "localhost"
  port: 5432
  database: "athema"
  username: "athema"
  password: "secret"
  sslmode: "disable"
log:
  level: "info"
  format: "json"
memory:
  enabled: true
conversation:
  enabled: true
`
	if err := os.WriteFile(cfgPath, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Override subsystem enabled via env var (auto-derived: ATHEMA_MEMORY_ENABLED).
	t.Setenv("ATHEMA_MEMORY_ENABLED", "false")
	t.Setenv("ATHEMA_CONVERSATION_ENABLED", "false")

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Memory.Enabled {
		t.Error("Memory.Enabled = true, want false (env override)")
	}
	if cfg.Conversation.Enabled {
		t.Error("Conversation.Enabled = true, want false (env override)")
	}
}

func TestLoad_InvalidEnvValue(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")

	yamlContent := `
server:
  host: "localhost"
  port: 8080
db:
  host: "localhost"
  port: 5432
  database: "athema"
  username: "athema"
  password: "secret"
  sslmode: "disable"
log:
  level: "info"
  format: "json"
`
	if err := os.WriteFile(cfgPath, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("ATHEMA_SERVER_PORT", "not_a_number")

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// Invalid env value should keep YAML default
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want %d (YAML default when env is invalid)", cfg.Server.Port, 8080)
	}
}
