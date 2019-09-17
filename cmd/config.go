package cmd

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Host             string `yaml:"url"`
	Port             int    `yaml:"port"`
	AutoScalingGroup string `yaml:"auto_scaling_group"`
}

type ErrFileNotExist struct {
	message string
}

func (e ErrFileNotExist) Error() string {
	return e.message
}

func GetConfig(configFile string) (*Config, error) {
	f, err := os.OpenFile(configFile, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &ErrFileNotExist{message: err.Error()}
		}
		return nil, err
	}

	buf, err := ioutil.ReadAll(f)
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
