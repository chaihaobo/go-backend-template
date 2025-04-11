package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type (
	// Configuration is the configuration of the application.
	Configuration struct {
		*viper.Viper
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
	vp := viper.New()
	vp.SetConfigFile(path)
	vp.SetConfigType("yaml")
	configuration := new(Configuration)
	configuration.Viper = vp
	remoteConfigSetupOK, err := configuration.setupRemoteConfig()
	if err != nil {
		return nil, err
	}
	if remoteConfigSetupOK {
		return configuration, nil
	}

	if err := configuration.setupLocalConfig(); err != nil {
		return nil, err
	}
	return configuration, nil
}

func (c *Configuration) setupRemoteConfig() (bool, error) {
	// 判断配置中心配置是否存在
	envConsulConfigEndpoint := os.Getenv("CONSUL_CONFIG_ENDPOINT")
	envConsulConfigPath := os.Getenv("CONSUL_CONFIG_PATH")
	if envConsulConfigPath == "" || envConsulConfigEndpoint == "" {
		return false, nil
	}
	if err := c.AddRemoteProvider("consul", envConsulConfigEndpoint, envConsulConfigPath); err != nil {
		return false, err
	}
	if err := c.ReadRemoteConfig(); err != nil {
		return false, err
	}
	if err := c.Unmarshal(c); err != nil {
		return false, err
	}
	go c.watchRemote()
	return true, nil

}

func (c *Configuration) setupLocalConfig() error {
	if err := c.ReadInConfig(); err != nil {
		return err
	}
	if err := c.Unmarshal(c); err != nil {
		return err
	}
	c.WatchConfig()
	return nil
}

func (c *Configuration) watchRemote() {
	for {
		select {
		case <-time.Tick(time.Second * 10):
			c.pullRemoteConfig()
		}

	}
}

func (c *Configuration) pullRemoteConfig() {
	if err := c.WatchRemoteConfig(); err != nil {
		slog.Error("watch remote config failed", slog.String("error", err.Error()))
		return
	}
	if err := c.Unmarshal(c); err != nil {
		slog.Error("unmarshal remote config failed", slog.String("error", err.Error()))
		return
	}
	return
}
