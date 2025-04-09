package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	// Configuration is the configuration of the application.
	Configuration struct {
		Service  Service  `yaml:"service"`
		Database Database `yaml:"database"`
		Logger   Logger   `yaml:"logger"`
		Redis    Redis    `yaml:"redis"`
		JWT      JWT      `yaml:"jwt"`
	}
	// Service is the configuration of the service.
	Service struct {
		Name               string `yaml:"name"`
		HTTPPort           int    `yaml:"httpPort"`
		GrpcPort           int    `yaml:"grpcPort"`
		MetricPort         int    `yaml:"metricPort"`
		Debug              bool   `yaml:"debug"`
		TraceCollectorURL  string `yaml:"traceCollectorURL"`
		Env                string `yaml:"env"`
		DiscoveryServerURL string `yaml:"discoveryServerURL"`
	}

	// Database is the configuration of the database.
	Database struct {
		Host        string        `yaml:"host"`
		Port        string        `yaml:"port"`
		User        string        `yaml:"user"`
		Password    string        `yaml:"password"`
		Name        string        `yaml:"name"`
		MaxOpen     int           `yaml:"maxOpen"`
		MaxIdle     int           `yaml:"maxIdle"`
		MaxLifetime time.Duration `yaml:"maxLifetime"`
		MaxIdleTime time.Duration `yaml:"maxIdleTime"`
		Location    string        `yaml:"location"`
	}

	// Logger is the configuration of the logger.
	Logger struct {
		FileName string
		MaxSize  int
		MaxAge   int
	}

	// Redis is the configuration of the redis.
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		Index    int    `yaml:"index"`
	}

	// JWT is the configuration of the jwt.
	JWT struct {
		AccessTokenSecretKey  string        `yaml:"accessTokenSecretKey"`
		RefreshTokenSecretKey string        `yaml:"refreshTokenSecretKey"`
		AccessTokenDuration   time.Duration `yaml:"accessTokenDuration"`
		RefreshTokenDuration  time.Duration `yaml:"refreshTokenDuration"`
	}
)

func NewConfiguration(path string) (*Configuration, error) {
	configRawFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configRawFile.Close()
	configuration := new(Configuration)

	yamlDecoder := yaml.NewDecoder(configRawFile)
	if err := yamlDecoder.Decode(configuration); err != nil {
		return nil, err
	}
	return configuration, nil
}
