package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Upload   UploadConfig   `mapstructure:"upload"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"` // debug or release
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Source string `mapstructure:"source"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	AccessExpiry  time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
}

type LogConfig struct {
	Level       string            `mapstructure:"level"`
	Format      string            `mapstructure:"format"`
	Output      LogOutputConfig   `mapstructure:"output"`
	Rotation    LogRotationConfig `mapstructure:"rotation"`
	Caller      bool              `mapstructure:"caller"`
	ServiceName string            `mapstructure:"service_name"`
}

type LogOutputConfig struct {
	ToStdout bool   `mapstructure:"to_stdout"`
	ToFile   bool   `mapstructure:"to_file"`
	FilePath string `mapstructure:"file_path"`
}

type LogRotationConfig struct {
	MaxSizeMB  int  `mapstructure:"max_size_mb"`
	MaxBackups int  `mapstructure:"max_backups"`
	MaxAgeDays int  `mapstructure:"max_age_days"`
	Compress   bool `mapstructure:"compress"`
}

type CORSConfig struct {
	Enabled          bool     `mapstructure:"enabled"`
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

type UploadConfig struct {
	Dir string `mapstructure:"dir"`
}

func (d *DatabaseConfig) DSN() string {
	return d.Source
}

func (s *ServerConfig) IsDebug() bool {
	return s.Mode == "debug"
}

// DataDir returns the data directory path.
// Priority: DATA_DIR env > /app/data (Docker default) > ./data
func DataDir() string {
	if dir := os.Getenv("DATA_DIR"); dir != "" {
		return dir
	}
	// Docker container: /app/data
	if _, err := os.Stat("/app"); err == nil {
		return "/app/data"
	}
	return "./data"
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Check if CONFIG_PATH is set (启动参数传入)
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// Auto search paths if no explicit path
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("/app")
		viper.AddConfigPath("/etc/oa-nsdiy")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(viper.GetViper())

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	// Register a comma-separated StringToSlice decode hook so that env vars
	// like CORS_ALLOW_ORIGINS=http://a,http://b are split into []string
	// correctly. Without this, AutomaticEnv returns the raw string and the
	// slice field ends up as a single-element ["a,b"]. The existing default
	// decode hook (which includes StringToTimeDurationHookFunc) is preserved
	// via ComposeDecodeHookFunc.
	if err := viper.Unmarshal(&cfg, func(dc *mapstructure.DecoderConfig) {
		dc.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			dc.DecodeHook,
			mapstructure.StringToSliceHookFunc(","),
		)
	}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Auto-derive database source path from DATA_DIR if using default
	if cfg.Database.Driver == "sqlite" && (cfg.Database.Source == "" || cfg.Database.Source == "./data/oa_nsdiy.db" || cfg.Database.Source == "./data/db/oa_nsdiy.db") {
		cfg.Database.Source = filepath.Join(DataDir(), "db", "oa_nsdiy.db")
	}

	// Auto-derive log file path from DATA_DIR if using default
	if cfg.Log.Output.ToFile && cfg.Log.Output.FilePath == "" {
		cfg.Log.Output.FilePath = filepath.Join(DataDir(), "logs", "oa-nsdiy.log")
	}

	// Auto-derive upload dir from DATA_DIR if using default
	if cfg.Upload.Dir == "" || cfg.Upload.Dir == "./data/uploads" {
		cfg.Upload.Dir = filepath.Join(DataDir(), "uploads")
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 3001)
	v.SetDefault("server.mode", "release")
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")

	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.source", "./data/db/oa_nsdiy.db")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("jwt.secret", "default-secret-change-me")
	v.SetDefault("jwt.access_expiry", "30m")
	v.SetDefault("jwt.refresh_expiry", "168h")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output.to_stdout", true)
	v.SetDefault("log.output.to_file", false)
	v.SetDefault("log.output.file_path", "")
	v.SetDefault("log.rotation.max_size_mb", 100)
	v.SetDefault("log.rotation.max_backups", 10)
	v.SetDefault("log.rotation.max_age_days", 7)
	v.SetDefault("log.rotation.compress", true)
	v.SetDefault("log.caller", false)
	v.SetDefault("log.service_name", "oa-nsdiy")

	v.SetDefault("cors.enabled", true)
	v.SetDefault("cors.allow_origins", []string{"http://localhost:3000"})
	v.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allow_headers", []string{"Content-Type", "Authorization"})
	v.SetDefault("cors.allow_credentials", true)

	v.SetDefault("upload.dir", "./data/uploads")
}
