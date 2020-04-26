package main

import (
	"chillit-vk-bot/internal/app/vkbot"
	"log"
)

func main() {
	if err := vkbot.Start(&vkbot.Config{
		GroupID:      0,
		Confirmation: "",
		Host:         ":8973",
	}); err != nil {
		log.Println(err)
	}
}
