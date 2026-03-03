package config

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
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
		url.QueryEscape(c.Username), url.QueryEscape(c.Password),
		c.Host, c.Port, c.Database, c.SSLMode,
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
	Enabled bool `yaml:"enabled"`
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
// Fields with explicit env tags use those names. Fields without env tags auto-derive
// the env var name from the config hierarchy (e.g., ATHEMA_MEMORY_ENABLED).
func applyEnvOverrides(cfg *Config) {
	applyEnvToStruct(reflect.ValueOf(cfg).Elem(), "ATHEMA")
}

func applyEnvToStruct(v reflect.Value, prefix string) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		if fieldVal.Kind() == reflect.Struct {
			yamlTag := field.Tag.Get("yaml")
			nestedPrefix := prefix
			if yamlTag != "" && yamlTag != "-" {
				nestedPrefix = prefix + "_" + strings.ToUpper(yamlTag)
			}
			applyEnvToStruct(fieldVal, nestedPrefix)
			continue
		}

		envTag := field.Tag.Get("env")
		if envTag == "" {
			// Auto-derive env var name from prefix + yaml tag.
			yamlTag := field.Tag.Get("yaml")
			if yamlTag != "" && yamlTag != "-" && prefix != "" {
				envTag = prefix + "_" + strings.ToUpper(yamlTag)
			}
		}
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
			intVal, err := strconv.Atoi(envVal)
			if err != nil {
				fmt.Fprintf(os.Stderr, "config: env %s=%q is not a valid integer, keeping YAML default\n", envTag, envVal)
				continue
			}
			fieldVal.SetInt(int64(intVal))
		case reflect.Bool:
			switch strings.ToLower(envVal) {
			case "true", "1":
				fieldVal.SetBool(true)
			case "false", "0":
				fieldVal.SetBool(false)
			default:
				fmt.Fprintf(os.Stderr, "config: env %s=%q is not a valid boolean, keeping YAML default\n", envTag, envVal)
			}
		}
	}
}
