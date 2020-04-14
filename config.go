package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	GRPC GRPCConfig `yaml:"grpc"`
	HTTP HTTPConfig `yaml:"http"`
}

type GRPCConfig struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

func NewConfig(path string) (Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	conf := Config{}
	if err := yaml.Unmarshal(file, &conf); err != nil {
		return Config{}, err
	}
	return conf, nil
}
