package conf

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server Server `yaml:"server"`
	Micro  Micro  `yaml:"micro"`
	Etcd   Etcd   `yaml:"etcd"`
}

type Server struct {
	Addr string `yaml:"addr"`
}

type Micro struct {
	Name string `yaml:"name"`
}

type Etcd struct {
	Addrs []string `yaml:"addrs"`
}

func LoadConfig(confPath string) (*Config, error) {
	config := &Config{}
	data, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
