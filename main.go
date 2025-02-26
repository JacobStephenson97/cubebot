package main

import (
	bot "cubebot/bot"
	"cubebot/internal/db"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	if db != nil {
		fmt.Println("DB connected!")
	}

	bot.BotToken = os.Getenv("DISCORD_TOKEN")
	bot.Run(db) // Pass the DB connection to Run
}
