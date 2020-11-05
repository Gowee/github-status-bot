package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gowee/github-status-bot/internal/bot"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not specified")
		return
	}

	chatID := os.Getenv("CHAT_ID")
	if chatID == "" {
		log.Fatal("CHAT_ID is unspecified")
		return
	}

	dataFilePath := os.Getenv("DATA_FILE_PATH")
	if chatID == "" {
		dataFilePath = "./data.json"
	}

	interval, _ := strconv.ParseInt(os.Getenv("CHECK_INTERVAL"), 10, 32)

	bot := bot.NewBotFromOptions(
		bot.Options{
			BotToken:      token,
			ChatID:        chatID,
			DataFilePath:  dataFilePath,
			CheckInterval: time.Duration(interval) * time.Second,
		},
	)
	bot.Run()
}
