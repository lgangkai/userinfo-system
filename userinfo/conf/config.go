package conf

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	MysqlMaster *Mysql `yaml:"mysql-master"`
	MysqlSlave  *Mysql `yaml:"mysql-slave"`
	Redis       *Redis `yaml:"redis"`
	Etcd        *Etcd  `yaml:"etcd"`
	Micro       *Micro `yaml:"micro"`
}

type Mysql struct {
	Driver   string `yaml:"driver"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       string `yaml:"db"`
}

type Redis struct {
	Addrs []string `yaml:"addrs"`
}

type Etcd struct {
	Addrs []string `yaml:"addrs"`
}

type Micro struct {
	Name string `yaml:"name"`
	Addr string `yaml:"addr"`
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
