package main

import (
	"chillit-vk-bot/internal/app/config"
	"chillit-vk-bot/internal/app/places"
	"chillit-vk-bot/internal/app/vkbot"
	"flag"
	"log"

	"google.golang.org/grpc"
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

	log.Printf("Connecting to places store service '%s'...\n", configuration.StoreService.URL)
	conn, err := grpc.Dial(
		configuration.StoreService.URL,
		grpc.WithInsecure(),
		// grpc.WithBlock(),
		// grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Fatalf("could not connect to store service. %v\n", err)
	}
	defer conn.Close()
	log.Println("Connected!")

	log.Println("Starting...")
	if err := vkbot.Start(configuration.VkBot, places.NewPlacesStoreClient(conn)); err != nil {
		log.Println(err)
	}
	log.Println("Shutted down")
}
