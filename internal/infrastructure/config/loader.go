package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	DB       DBConfig       `yaml:"db"`
	LLM      LLMConfig      `yaml:"llm"`
	Log      LogConfig      `yaml:"log"`
	Memory   SubsystemConfig `yaml:"memory"`
	Conversation SubsystemConfig `yaml:"conversation"`
	Personality  SubsystemConfig `yaml:"personality"`
	Emotional    SubsystemConfig `yaml:"emotional"`
	Lifecycle    SubsystemConfig `yaml:"lifecycle"`
	Initiation   SubsystemConfig `yaml:"initiation"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host string `yaml:"host" env:"ATHEMA_SERVER_HOST"`
	Port int    `yaml:"port" env:"ATHEMA_SERVER_PORT"`
}

// DBConfig holds database connection settings.
type DBConfig struct {
	Host     string `yaml:"host"     env:"ATHEMA_DB_HOST"`
	Port     int    `yaml:"port"     env:"ATHEMA_DB_PORT"`
	Database string `yaml:"database" env:"ATHEMA_DB_DATABASE"`
	Username string `yaml:"username" env:"ATHEMA_DB_USERNAME"`
	Password string `yaml:"password" env:"ATHEMA_DB_PASSWORD"`
	SSLMode  string `yaml:"sslmode"  env:"ATHEMA_DB_SSLMODE"`
}

// DSN returns the PostgreSQL connection string.
func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
	)
}

// LLMConfig holds LLM provider settings.
type LLMConfig struct {
	Provider string `yaml:"provider" env:"ATHEMA_LLM_PROVIDER"`
	APIKey   string `yaml:"api_key"  env:"ATHEMA_LLM_API_KEY"`
	Model    string `yaml:"model"    env:"ATHEMA_LLM_MODEL"`
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level  string `yaml:"level"  env:"ATHEMA_LOG_LEVEL"`
	Format string `yaml:"format" env:"ATHEMA_LOG_FORMAT"`
}

// SubsystemConfig holds per-subsystem configuration stubs.
type SubsystemConfig struct {
	Enabled bool `yaml:"enabled" env:""`
}

// Load reads configuration from a YAML file and applies environment variable overrides.
// Environment variables take precedence over YAML values.
func Load(path string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %s: %w", path, err)
	}

	applyEnvOverrides(cfg)

	return cfg, nil
}

// applyEnvOverrides walks the config struct and overrides values with matching environment variables.
func applyEnvOverrides(cfg *Config) {
	applyEnvToStruct(reflect.ValueOf(cfg).Elem())
}

func applyEnvToStruct(v reflect.Value) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		if fieldVal.Kind() == reflect.Struct {
			applyEnvToStruct(fieldVal)
			continue
		}

		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envVal, ok := os.LookupEnv(envTag)
		if !ok {
			continue
		}

		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(envVal)
		case reflect.Int:
			var intVal int
			if _, err := fmt.Sscanf(envVal, "%d", &intVal); err == nil {
				fieldVal.SetInt(int64(intVal))
			}
		case reflect.Bool:
			fieldVal.SetBool(strings.EqualFold(envVal, "true") || envVal == "1")
		}
	}
}
