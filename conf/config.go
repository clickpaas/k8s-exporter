package conf

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

var (
	cfg  *Config
	once sync.Once
)

type Config struct {
	Cluster string `json:"cluster"`
	MonitorDomain []struct{
		Name string `json:"name"`
		Domain string `json:"domain"`
	} `json:"monitorDomain"`
}

func NewConfig(cfgFile string) *Config {
	viper.SetConfigFile(cfgFile)
	once.Do(func() {
		cfg = &Config{}
		if err := viper.ReadInConfig(); err != nil {
			logrus.Fatalf("read config failed,reason: %s", err.Error())
		}
		if err := viper.Unmarshal(cfg); err != nil {
			logrus.Fatal("unmarshal config failed, reason: %s", err.Error())
		}
	})
	return cfg
}
