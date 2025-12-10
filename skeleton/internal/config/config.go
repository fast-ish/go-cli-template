// Package config handles application configuration
package config

import (
	"fmt"
	"os"
	"path/filepath"

{%- if values.cliFramework == "cobra" %}
	"github.com/spf13/viper"
{%- endif %}
{%- if values.configFormat == "yaml" or values.configFormat == "all" %}
	"gopkg.in/yaml.v3"
{%- endif %}
{%- if values.configFormat == "toml" or values.configFormat == "all" %}
	"github.com/pelletier/go-toml/v2"
{%- endif %}
)

// Config represents the application configuration
type Config struct {
{%- if values.aiProvider != "none" %}
	AI AIConfig `json:"ai" yaml:"ai" toml:"ai"`
{%- endif %}
{%- for integration in values.integrations %}
	{{integration|title}} {{integration|title}}Config `json:"{{integration}}" yaml:"{{integration}}" toml:"{{integration}}"`
{%- endfor %}
	Logging LoggingConfig `json:"logging" yaml:"logging" toml:"logging"`
{%- if values.metrics %}
	Metrics MetricsConfig `json:"metrics" yaml:"metrics" toml:"metrics"`
{%- endif %}
{%- if values.tracing %}
	Tracing TracingConfig `json:"tracing" yaml:"tracing" toml:"tracing"`
{%- endif %}
}

{%- if values.aiProvider != "none" %}

// AIConfig holds AI provider configuration
type AIConfig struct {
	Provider string `json:"provider" yaml:"provider" toml:"provider"`
{%- if values.aiProvider == "bedrock" %}
	Region   string `json:"region" yaml:"region" toml:"region"`
	Model    string `json:"model" yaml:"model" toml:"model"`
{%- elif values.aiProvider == "openai" %}
	APIKey   string `json:"api_key" yaml:"api_key" toml:"api_key"`
	Model    string `json:"model" yaml:"model" toml:"model"`
{%- elif values.aiProvider == "anthropic" %}
	APIKey   string `json:"api_key" yaml:"api_key" toml:"api_key"`
	Model    string `json:"model" yaml:"model" toml:"model"`
{%- elif values.aiProvider == "ollama" %}
	Host     string `json:"host" yaml:"host" toml:"host"`
	Model    string `json:"model" yaml:"model" toml:"model"`
{%- endif %}
}
{%- endif %}

{%- for integration in values.integrations %}

// {{integration|title}}Config holds {{integration}} configuration
type {{integration|title}}Config struct {
{%- if integration == "aws" %}
	Region  string `json:"region" yaml:"region" toml:"region"`
	Profile string `json:"profile" yaml:"profile" toml:"profile"`
{%- elif integration == "github" %}
	Token string `json:"token" yaml:"token" toml:"token"`
	Org   string `json:"org" yaml:"org" toml:"org"`
{%- elif integration == "kubernetes" %}
	Kubeconfig string `json:"kubeconfig" yaml:"kubeconfig" toml:"kubeconfig"`
	Context    string `json:"context" yaml:"context" toml:"context"`
{%- elif integration == "grafana" %}
	URL   string `json:"url" yaml:"url" toml:"url"`
	Token string `json:"token" yaml:"token" toml:"token"`
{%- elif integration == "slack" %}
	Token   string `json:"token" yaml:"token" toml:"token"`
	Channel string `json:"channel" yaml:"channel" toml:"channel"`
{%- elif integration == "notion" %}
	Token string `json:"token" yaml:"token" toml:"token"`
{%- elif integration == "argocd" %}
	URL      string `json:"url" yaml:"url" toml:"url"`
	Token    string `json:"token" yaml:"token" toml:"token"`
	Insecure bool   `json:"insecure" yaml:"insecure" toml:"insecure"`
{%- elif integration == "terraform" %}
	WorkingDir string `json:"working_dir" yaml:"working_dir" toml:"working_dir"`
{%- endif %}
}
{%- endfor %}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `json:"level" yaml:"level" toml:"level"`
	Format string `json:"format" yaml:"format" toml:"format"` // json or text
}

{%- if values.metrics %}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Enabled bool   `json:"enabled" yaml:"enabled" toml:"enabled"`
	Port    int    `json:"port" yaml:"port" toml:"port"`
	Path    string `json:"path" yaml:"path" toml:"path"`
}
{%- endif %}

{%- if values.tracing %}

// TracingConfig holds tracing configuration
type TracingConfig struct {
	Enabled  bool   `json:"enabled" yaml:"enabled" toml:"enabled"`
	Endpoint string `json:"endpoint" yaml:"endpoint" toml:"endpoint"`
	Service  string `json:"service" yaml:"service" toml:"service"`
}
{%- endif %}

// Load loads configuration from file
func Load(configFile string) (*Config, error) {
{%- if values.cliFramework == "cobra" %}
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Config file paths to try
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}

		v.SetConfigName("config")
{%- if values.configFormat == "yaml" %}
		v.SetConfigType("yaml")
{%- elif values.configFormat == "toml" %}
		v.SetConfigType("toml")
{%- elif values.configFormat == "json" %}
		v.SetConfigType("json")
{%- elif values.configFormat == "all" %}
		v.SetConfigType("yaml") // default
{%- endif %}
		v.AddConfigPath(filepath.Join(home, ".${{values.name}}"))
		v.AddConfigPath(".")
	}

	// Environment variables
	v.SetEnvPrefix("{{values.name|upper}}")
	v.AutomaticEnv()

	// Read config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// Config file not found, use defaults
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
{%- else %}
	// Simple config loading without viper
	if configFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configFile = filepath.Join(home, ".${{values.name}}", "config.yaml")
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Use defaults
			return defaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
{%- if values.configFormat == "yaml" %}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
{%- elif values.configFormat == "toml" %}
	if err := toml.Unmarshal(data, &cfg); err != nil {
{%- elif values.configFormat == "json" %}
	if err := json.Unmarshal(data, &cfg); err != nil {
{%- elif values.configFormat == "all" %}
	// Try YAML first
	if err := yaml.Unmarshal(data, &cfg); err != nil {
{%- endif %}
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
{%- endif %}
}

{%- if values.cliFramework == "cobra" %}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "text")
{%- if values.metrics %}
	v.SetDefault("metrics.enabled", false)
	v.SetDefault("metrics.port", 9090)
	v.SetDefault("metrics.path", "/metrics")
{%- endif %}
{%- if values.tracing %}
	v.SetDefault("tracing.enabled", false)
	v.SetDefault("tracing.endpoint", "http://localhost:4318")
	v.SetDefault("tracing.service", "${{values.name}}")
{%- endif %}
{%- if values.aiProvider != "none" %}
	v.SetDefault("ai.provider", "${{values.aiProvider}}")
{%- if values.aiProvider == "bedrock" %}
	v.SetDefault("ai.region", "us-west-2")
	v.SetDefault("ai.model", "anthropic.claude-3-sonnet-20240229-v1:0")
{%- elif values.aiProvider == "openai" %}
	v.SetDefault("ai.model", "gpt-4")
{%- elif values.aiProvider == "anthropic" %}
	v.SetDefault("ai.model", "claude-3-5-sonnet-20241022")
{%- elif values.aiProvider == "ollama" %}
	v.SetDefault("ai.host", "http://localhost:11434")
	v.SetDefault("ai.model", "llama2")
{%- endif %}
{%- endif %}
}
{%- else %}

// defaultConfig returns default configuration
func defaultConfig() *Config {
	return &Config{
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
{%- if values.metrics %}
		Metrics: MetricsConfig{
			Enabled: false,
			Port:    9090,
			Path:    "/metrics",
		},
{%- endif %}
{%- if values.tracing %}
		Tracing: TracingConfig{
			Enabled:  false,
			Endpoint: "http://localhost:4318",
			Service:  "${{values.name}}",
		},
{%- endif %}
{%- if values.aiProvider != "none" %}
		AI: AIConfig{
			Provider: "${{values.aiProvider}}",
{%- if values.aiProvider == "bedrock" %}
			Region:   "us-west-2",
			Model:    "anthropic.claude-3-sonnet-20240229-v1:0",
{%- elif values.aiProvider == "openai" %}
			Model:    "gpt-4",
{%- elif values.aiProvider == "anthropic" %}
			Model:    "claude-3-5-sonnet-20241022",
{%- elif values.aiProvider == "ollama" %}
			Host:     "http://localhost:11434",
			Model:    "llama2",
{%- endif %}
		},
{%- endif %}
	}
}
{%- endif %}
