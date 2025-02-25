package main

import (
	bot "cubebot/bot"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot.BotToken = os.Getenv("DISCORD_TOKEN")

	bot.Run() // call the run function of bot/bot.go
}
