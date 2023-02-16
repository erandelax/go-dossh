package configuration

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var cfg AppConfig

type AppConfig struct {
	Host       string                     `yaml:"host"`
	Port       int                        `yaml:"port"`
	Users      map[string]UserConfig      `yaml:"users"`
	Containers map[string]ContainerConfig `yaml:"containers"`
}

type UserConfig struct {
	PublicKey  string `yaml:"public_key"`
	Containers map[string][]string
}

type ContainerConfig struct {
	As *string `yaml:"as"`
}

func init() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Panicln(err)
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func Get() AppConfig {
	return cfg
}
