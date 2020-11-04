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

	chatID, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if chatID == 0 || err != nil {
		log.Fatal("CHAT_ID is unspecified or invalid")
		return
	}

	interval, _ := strconv.ParseInt(os.Getenv("CHECK_INTERVAL"), 10, 32)
	// log.Println("Check interval: ", interval)

	bot := bot.NewBotFromOptions(bot.Options{BotToken: token, ChatID: chatID, DataFilePath: "./data.json", CheckInterval: time.Duration(interval) * time.Second})
	bot.Run()
}
