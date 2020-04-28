package main

import (
	"chillit-vk-bot/internal/app/config"
	"chillit-vk-bot/internal/app/vkbot"
	"flag"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config_path", "configs/config.yaml", "path to config file")
}

func main() {
	flag.Parse()

	configuration, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalln("error loading configuration", err)
	}

	log.Println("Started!")
	if err := vkbot.Start(configuration.VkBot); err != nil {
		log.Println(err)
	}
	log.Println("Shutted down...")
}
