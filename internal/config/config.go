package config

import "fmt"

type Config struct {
	LogLevel     string   `hcl:"log_level,optional"`
	PoolName     string   `hcl:"pool_name,optional"`
	PluginDir    string   `hcl:"plugin_dir,optional"`
	EmbedPlugins []string `hcl:"embed_plugins,optional"`

	Server  *ServerConfig  `hcl:"server,block"`
	Client  *ClientConfig  `hcl:"client,block"`
	Metrics *MetricsConfig `hcl:"metrics,block"`
	Redis   *RedisConfig   `hcl:"redis,block"`

	Plugins []*PluginConfig `hcl:"plugin,block"`
}

func CreateDefault() *Config {
	return &Config{
		LogLevel:     "INFO",
		PoolName:     "default",
		PluginDir:    "./plugins",
		EmbedPlugins: make([]string, 0),

		Server: &ServerConfig{
			Enabled: true,
			Address: "0.0.0.0:8080",
		},
		Client: &ClientConfig{
			Enabled: true,
		},
		Metrics: &MetricsConfig{
			Enabled: true,
			Address: "127.0.0.1:9001",
		},
		Redis: &RedisConfig{
			Endpoint: "127.0.0.1:6379",
			Database: 0,
			Password: "",
		},

		Plugins: make([]*PluginConfig, 0),
	}
}

func (c *Config) ValidateConfig() error {
	if c.LogLevel == "" {
		return fmt.Errorf("'log_level' is required")
	}
	if c.PluginDir == "" {
		return fmt.Errorf("'plugin_dir' is required")
	}

	if c.Server == nil {
		return fmt.Errorf("block 'server' is required")
	}
	if err := c.Server.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid 'server' block: %w", err)
	}

	if c.Client == nil {
		return fmt.Errorf("block 'client' is required")
	}
	if err := c.Client.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid 'client' block: %w", err)
	}

	if c.Metrics == nil {
		return fmt.Errorf("block 'metrics' is required")
	}
	if err := c.Metrics.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid 'metrics' block: %w", err)
	}

	if c.Redis == nil {
		return fmt.Errorf("block 'redis' is required")
	}
	if err := c.Redis.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid 'redis' block: %w", err)
	}

	return nil
}
