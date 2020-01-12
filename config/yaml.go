package config

import (
	"gitlab.com/pangold/auth/utils"
	"gopkg.in/yaml.v2"
)

type Yaml struct {
	path string
}

func NewYaml(path string) *Yaml {
	return &Yaml{
		path: path,
	}
}

func (this *Yaml) ReadConfig() (*Config) {
	conf := Config{}
	data, err := utils.ReadFile(this.path)
	if err != nil {
		panic(err.Error())
	}
	if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
		panic(err.Error())
	}
	return &conf
}

func (this *Yaml) WriteConfig(conf Config) error {
	data, err := yaml.Marshal(&conf)
	if err != nil {
		return err
	}
	if err := utils.WriteFile(this.path, string(data)); err != nil {
		return err
	}
	return nil
}
