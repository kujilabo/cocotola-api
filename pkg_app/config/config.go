package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

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

type DebugConfig struct {
	GinMode bool `yaml:"ginMode"`
	Wait    bool `yaml:"wait"`
}

type Config struct {
	CORS     *CORSConfig     `yaml:"cors"`
	Shutdown *ShutdownConfig `yaml:"shutdown"`
	Log      *LogConfig      `yaml:"log"`
	Debug    *DebugConfig    `yaml:"debug"`
}

func LoadConfig(env string) (*Config, error) {
	confContent, err := ioutil.ReadFile("./configs/" + env + ".yml")
	if err != nil {
		return nil, err
	}

	confContent = []byte(os.ExpandEnv(string(confContent)))
	conf := &Config{}
	if err := yaml.Unmarshal(confContent, conf); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(conf); err != nil {
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
