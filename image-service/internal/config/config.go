package config

import (
	"flag"
	"github.com/d1rtyloudx/spotiby-pkg/minio"
	"github.com/d1rtyloudx/spotiby-pkg/rabbitmq"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	HTTP     HTTPConfig     `yaml:"http"`
	Minio    MinioConfig    `yaml:"minio"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
}

type HTTPConfig struct {
	Port            int           `yaml:"port"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type RabbitMQConfig struct {
	Connection rabbitmq.Config  `yaml:"connection"`
	Publishers PublishersConfig `yaml:"publishers"`
}

type PublishersConfig struct {
	ProfileImagePublisher  rabbitmq.PublisherConfig `yaml:"profile_image_publisher"`
	PlaylistImagePublisher rabbitmq.PublisherConfig `yaml:"playlist_image_publisher"`
	TrackImagePublisher    rabbitmq.PublisherConfig `yaml:"track_image_publisher"`
}

type MinioConfig struct {
	Connection minio.Config       `yaml:"connection"`
	Buckets    MinioBucketsConfig `yaml:"buckets"`
}

type MinioBucketsConfig struct {
	ProfileBucket  string `yaml:"profile_bucket"`
	TrackBucket    string `yaml:"track_bucket"`
	PlaylistBucket string `yaml:"playlist_bucket"`
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
