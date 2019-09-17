package cmd

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Host             string `yaml:"url"`
	Port             int    `yaml:"port"`
	AutoScalingGroup string `yaml:"auto_scaling_group"`
}

func GetConfig(configFile string) (*Config, error) {
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
