package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Service struct {
		Name              string `yaml:"name"`
		HTTPPort          string `yaml:"httpPort"`
		MetricPort        int    `yaml:"metricPort"`
		Debug             bool   `yaml:"debug"`
		TraceCollectorURL string `yaml:"traceCollectorURL"`
		Env               string `yaml:"env"`
	} `yaml:"service"`
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
	} `yaml:"database"`
	Logger struct {
		FileName string
		MaxSize  int
		MaxAge   int
	}
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		Index    int    `yaml:"index"`
	} `yaml:"redis"`
	JWT struct {
		AccessTokenSecretKey  string        `yaml:"accessTokenSecretKey"`
		RefreshTokenSecretKey string        `yaml:"refreshTokenSecretKey"`
		AccessTokenDuration   time.Duration `yaml:"accessTokenDuration"`
		RefreshTokenDuration  time.Duration `yaml:"refreshTokenDuration"`
	} `yaml:"jwt"`
}

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
