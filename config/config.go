package config

var cfg Config

func Get() *Config {
	return &cfg
}

type ServerConfig struct {
	Port int64 `yaml:"port"`
}

type LogConfig struct {
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
	Level      string `yaml:"level"`
}

type Config struct {
	ServerConfig ServerConfig             `yaml:"server"`
	LogConfig    LogConfig                `yaml:"log"`
	MCPConfig    map[string]ServiceConfig `yaml:"mcp_config"`
}

type ServiceConfig struct {
	BaseURL     string            `yaml:"base_url"`
	ExtraHeader map[string]string `yaml:"extra_header"`
	ConfigPath  string            `yaml:"config_path"`
}
