package config

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

type AppConfig struct {
	Name          string `yaml:"name" validate:"required"`
	Port          int    `yaml:"port" validate:"required"`
	OwnerPassword string `yaml:"ownerPassword" validate:"required"`
	TestUserEmail string `yaml:"testUserEmail" validate:"required"`
}

type SQLite3Config struct {
	File string `yaml:"file" validate:"required"`
}

type MySQLConfig struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	Database string `yaml:"database" validate:"required"`
}

type DBConfig struct {
	DriverName string         `yaml:"driverName"`
	SQLite3    *SQLite3Config `yaml:"sqlite3"`
	MySQL      *MySQLConfig   `yaml:"mysql"`
}

type AuthConfig struct {
	SigningKey          string `yaml:"signingKey"`
	AccessTokenTTLMin   int    `yaml:"accessTokenTtlMin" validate:"gte=1"`
	RefreshTokenTTLHour int    `yaml:"refreshTokenTtlHour" validate:"gte=1"`
	GoogleCallbackURL   string `yaml:"googleCallbackUrl" validate:"required"`
	GoogleClientID      string `yaml:"googleClientId" validate:"required"`
	GoogleClientSecret  string `yaml:"googleClientSecret" validate:"required"`
	APITimeoutSec       int    `yaml:"apiTimeoutSec" validate:"gte=1"`
}

type TranslatorConfig struct {
	Endpoint   string `yaml:"endpoint" validate:"required"`
	TimeoutSec int    `yaml:"timeoutSec" validate:"gte=1"`
	Username   string `yaml:"username" validate:"required"`
	Password   string `yaml:"password" validate:"required"`
}

type TatoebaConfig struct {
	Endpoint   string `yaml:"endpoint" validate:"required"`
	TimeoutSec int    `yaml:"timeoutSec" validate:"gte=1"`
	Username   string `yaml:"username" validate:"required"`
	Password   string `yaml:"password" validate:"required"`
}

type SynthesizerConfig struct {
	Endpoint   string `yaml:"endpoint" validate:"required"`
	TimeoutSec int    `yaml:"timeoutSec" validate:"gte=1"`
	Username   string `yaml:"username" validate:"required"`
	Password   string `yaml:"password" validate:"required"`
}

type JaegerConfig struct {
	Endpoint string `yaml:"endpoint" validate:"required"`
}

type TraceConfog struct {
	Exporter string        `yaml:"exporter" validate:"required"`
	Jaeger   *JaegerConfig `yaml:"jaeger"`
}

type CORSConfig struct {
	AllowOrigins []string `yaml:"allowOrigins"`
}

type ShutdownConfig struct {
	TimeSec1 int `yaml:"timeSec1" validate:"gte=1"`
	TimeSec2 int `yaml:"timeSec2" validate:"gte=1"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type SwaggerConfig struct {
	Enabled bool   `yaml:"enabled"`
	Host    string `yaml:"host"`
	Schema  string `yaml:"schema"`
}

type DebugConfig struct {
	GinMode bool `yaml:"ginMode"`
	Wait    bool `yaml:"wait"`
}

type Config struct {
	App         *AppConfig         `yaml:"app" validate:"required"`
	DB          *DBConfig          `yaml:"db" validate:"required"`
	Auth        *AuthConfig        `yaml:"auth" validate:"required"`
	Translator  *TranslatorConfig  `yaml:"translator" validate:"required"`
	Tatoeba     *TatoebaConfig     `yaml:"tatoeba" validate:"required"`
	Synthesizer *SynthesizerConfig `yaml:"synthesizer" validate:"required"`
	Trace       *TraceConfog       `yaml:"trace" validate:"required"`
	CORS        *CORSConfig        `yaml:"cors" validate:"required"`
	Shutdown    *ShutdownConfig    `yaml:"shutdown" validate:"required"`
	Log         *LogConfig         `yaml:"log" validate:"required"`
	Swagger     *SwaggerConfig     `yaml:"swagger" validate:"required"`
	Debug       *DebugConfig       `yaml:"debug"`
}

func LoadConfig(env string) (*Config, error) {
	confContent, err := os.ReadFile("./configs/" + env + ".yml")
	if err != nil {
		return nil, err
	}

	confContent = []byte(os.ExpandEnv(string(confContent)))
	conf := &Config{}
	if err := yaml.Unmarshal(confContent, conf); err != nil {
		return nil, err
	}

	if err := libD.Validator.Struct(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func InitLog(env string, cfg *LogConfig) error {
	formatter := &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	}
	logrus.SetFormatter(formatter)

	switch strings.ToLower(cfg.Level) {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.Infof("Unsupported log level: %s", cfg.Level)
		logrus.SetLevel(logrus.WarnLevel)
	}

	logrus.SetOutput(os.Stdout)

	return nil
}

func InitCORS(cfg *CORSConfig) cors.Config {
	if len(cfg.AllowOrigins) == 1 && cfg.AllowOrigins[0] == "*" {
		return cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"*"},
			AllowHeaders:    []string{"*"},
		}
	}

	return cors.Config{
		AllowOrigins: cfg.AllowOrigins,
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}
}
