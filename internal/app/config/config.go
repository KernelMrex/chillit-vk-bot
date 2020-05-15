package config

import (
	"chillit-vk-bot/internal/app/places"
	"chillit-vk-bot/internal/app/templates"
	"chillit-vk-bot/internal/app/vkbot"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Configuration structure which provides app configuration
type Configuration struct {
	VkBot        *vkbot.Config     `yaml:"vk_bot"`
	StoreService *places.Config    `yaml:"store_service"`
	Templates    *templates.Config `yaml:"templates"`
}

// NewConfig parse Configuration from yaml file
func NewConfig(path string) (*Configuration, error) {
	if !fileExists(path) {
		return nil, fmt.Errorf("[ NewConfig ] could not open file: file '%s' does not exists", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.New("[ NewConfig ] could not open file: " + err.Error())
	}
	configData, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.New("[ NewConfig ] error while reading file: " + err.Error())
	}

	config := &Configuration{}

	if err := yaml.Unmarshal(configData, config); err != nil {
		return nil, errors.New("[ NewConfig ] error while parsing config: " + err.Error())
	}

	return config, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
