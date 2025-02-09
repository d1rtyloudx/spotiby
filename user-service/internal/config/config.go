package config

import (
	"flag"
	"github.com/d1rtyloudx/spotiby-pkg/postgres"
	"github.com/d1rtyloudx/spotiby-pkg/rabbitmq"
	"github.com/d1rtyloudx/spotiby-pkg/redis"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	HTTP     HTTPConfig      `yaml:"http"`
	Postgres postgres.Config `yaml:"postgres"`
	RabbitMQ RabbitMQConfig  `yaml:"rabbitmq"`
	Token    TokenConfig     `yaml:"token"`
	Redis    redis.Config    `yaml:"redis"`
}

type HTTPConfig struct {
	Port            int           `yaml:"port"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type TokenConfig struct {
	AccessTTL     time.Duration `yaml:"access_ttl"`
	RefreshTTL    time.Duration `yaml:"refresh_ttl"`
	AccessSecret  string        `env:"ACCESS_SECRET"`
	RefreshSecret string        `env:"REFRESH_SECRET"`
}

type RabbitMQConfig struct {
	Connection   rabbitmq.Config                  `yaml:"connection"`
	ImageBinding rabbitmq.ExchangeAndQueueBinding `yaml:"image_binding"`
}

func MustLoad() *Config {
	cfgPath := fetchConfigPath()

	if cfgPath == "" {
		panic("config path not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		panic("config path does not exist: " + cfgPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
