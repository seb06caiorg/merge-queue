package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config represents the application configuration.
type Config struct {
	Server   ServerConfig   `json:"server"`
	App      AppConfig      `json:"app"`
	Features FeaturesConfig `json:"features"`
	Defaults DefaultsConfig `json:"defaults"`
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Port         string        `json:"port"`
	Host         string        `json:"host"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// AppConfig holds application-level configuration.
type AppConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Debug       bool   `json:"debug"`
	Environment string `json:"environment"` // "development", "staging", "production"
}

// FeaturesConfig holds feature flags and limits.
type FeaturesConfig struct {
	EnableCORS       bool `json:"enable_cors"`
	EnableLogging    bool `json:"enable_logging"`
	EnableMetrics    bool `json:"enable_metrics"`
	MaxTasksPerUser  int  `json:"max_tasks_per_user"`
	RateLimitPerMin  int  `json:"rate_limit_per_min"`
	EnableValidation bool `json:"enable_validation"`
}

// DefaultsConfig holds default values for various entities.
type DefaultsConfig struct {
	TaskStatus   string `json:"task_status"`
	TaskPriority string `json:"task_priority"`
	UserRole     string `json:"user_role"`
	PageSize     int    `json:"page_size"`
}

// LoadConfig loads configuration from a JSON file with environment variable overrides.
func LoadConfig(filename string) (*Config, error) {
	config := &Config{}

	// Set defaults first.
	config.setDefaults()

	// Load from file if it exists.
	if filename != "" {
		if err := config.loadFromFile(filename); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with environment variables.
	config.loadFromEnv()

	// Validate configuration.
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// setDefaults sets default configuration values.
func (c *Config) setDefaults() {
	c.Server = ServerConfig{
		Port:         ":8080",
		Host:         "localhost",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	c.App = AppConfig{
		Name:        "Task Manager API",
		Version:     "1.0.0",
		Debug:       false,
		Environment: "development",
	}

	c.Features = FeaturesConfig{
		EnableCORS:       true,
		EnableLogging:    true,
		EnableMetrics:    false,
		MaxTasksPerUser:  100,
		RateLimitPerMin:  60,
		EnableValidation: true,
	}

	c.Defaults = DefaultsConfig{
		TaskStatus:   "pending",
		TaskPriority: "medium",
		UserRole:     "user",
		PageSize:     20,
	}
}

// loadFromFile loads configuration from a JSON file.
func (c *Config) loadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		// File doesn't exist is not an error - we'll use defaults.
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(c)
}

// loadFromEnv loads configuration from environment variables.
func (c *Config) loadFromEnv() {
	if port := os.Getenv("PORT"); port != "" {
		if port[0] != ':' {
			port = ":" + port
		}
		c.Server.Port = port
	}

	if host := os.Getenv("HOST"); host != "" {
		c.Server.Host = host
	}

	if debug := os.Getenv("DEBUG"); debug != "" {
		c.App.Debug = debug == "true" || debug == "1"
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		c.App.Environment = env
	}

	if maxTasks := os.Getenv("MAX_TASKS_PER_USER"); maxTasks != "" {
		if val, err := strconv.Atoi(maxTasks); err == nil {
			c.Features.MaxTasksPerUser = val
		}
	}

	if rateLimit := os.Getenv("RATE_LIMIT_PER_MIN"); rateLimit != "" {
		if val, err := strconv.Atoi(rateLimit); err == nil {
			c.Features.RateLimitPerMin = val
		}
	}
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.App.Name == "" {
		return fmt.Errorf("app name is required")
	}

	if c.App.Version == "" {
		return fmt.Errorf("app version is required")
	}

	validEnvs := []string{"development", "staging", "production"}
	validEnv := false
	for _, env := range validEnvs {
		if c.App.Environment == env {
			validEnv = true
			break
		}
	}
	if !validEnv {
		return fmt.Errorf("invalid environment: %s", c.App.Environment)
	}

	if c.Features.MaxTasksPerUser <= 0 {
		return fmt.Errorf("max_tasks_per_user must be positive")
	}

	if c.Features.RateLimitPerMin <= 0 {
		return fmt.Errorf("rate_limit_per_min must be positive")
	}

	if c.Defaults.PageSize <= 0 {
		return fmt.Errorf("default page_size must be positive")
	}

	return nil
}

// IsDevelopment returns true if running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if running in production mode.
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// GetAddress returns the full server address.
func (c *Config) GetAddress() string {
	return c.Server.Host + c.Server.Port
}
