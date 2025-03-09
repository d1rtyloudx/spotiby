package config

import (
	"flag"
	"github.com/d1rtyloudx/spotiby-pkg/elastic"
	"github.com/d1rtyloudx/spotiby-pkg/kafka"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	HTTP    HTTPConfig    `yaml:"http"`
	Kafka   KafkaConfig   `yaml:"kafka"`
	Elastic ElasticConfig `yaml:"elastic"`
}

type ElasticConfig struct {
	Connection elastic.Config `yaml:"connection"`
	Indexes    ElasticIndexes `yaml:"indexes"`
}

type ElasticIndexes struct {
	ProfileIndex  elastic.IndexConfig `yaml:"profile_index"`
	PlaylistIndex elastic.IndexConfig `yaml:"playlist_index"`
	TrackIndex    elastic.IndexConfig `yaml:"track_index"`
}

type HTTPConfig struct {
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type KafkaConfig struct {
	Connection kafka.Config `yaml:"connection"`
	Topics     KafkaTopics  `yaml:"topics"`
}

type KafkaTopics struct {
	CreateProfileTopic  kafka.TopicConfig `yaml:"create_profile_topic"`
	UpdateProfileTopic  kafka.TopicConfig `yaml:"update_profile_topic"`
	DeleteProfileTopic  kafka.TopicConfig `yaml:"delete_profile_topic"`
	CreateTrackTopic    kafka.TopicConfig `yaml:"create_track_topic"`
	UpdateTrackTopic    kafka.TopicConfig `yaml:"update_track_topic"`
	DeleteTrackTopic    kafka.TopicConfig `yaml:"delete_track_topic"`
	CreatePlaylistTopic kafka.TopicConfig `yaml:"create_playlist_topic"`
	UpdatePlaylistTopic kafka.TopicConfig `yaml:"update_playlist_topic"`
	DeletePlaylistTopic kafka.TopicConfig `yaml:"delete_playlist_topic"`
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
