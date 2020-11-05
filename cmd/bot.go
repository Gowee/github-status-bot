package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gowee/github-status-bot/internal/bot"
	"github.com/gowee/github-status-bot/pkg/utils"
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
	if dataFilePath == "" {
		dataFilePath = "./data.json"
	}

	interval, _ := strconv.ParseInt(os.Getenv("CHECK_INTERVAL"), 10, 32)

	chatDescriptionTemplate := strings.TrimSpace(os.Getenv("CHAT_DESCRIPTION_TEMPLATE"))
	if chatDescriptionTemplate != "" && !strings.Contains(chatDescriptionTemplate, "%s") {
		var err error
		// WTF: using := here will shadow the outer chatDescriptionTemplate.
		//      it seems to be a design flaw in the overall syntax.
		chatDescriptionTemplate, err = utils.B64Dec(chatDescriptionTemplate)
		if err != nil || !strings.Contains(chatDescriptionTemplate, "%s") {
			log.Fatal(
				"CHAT_DESCRIPTION_TEMPLATE expects a \"%s\" which will be replaced to generated content",
			)
			return
		}
		// chatDescriptionTemplate = "%s\n\n\nPowered by https://git.io/ghstsbot"
	}

	bot := bot.NewBotFromOptions(
		bot.Options{
			BotToken:                token,
			ChatID:                  chatID,
			DataFilePath:            dataFilePath,
			CheckInterval:           time.Duration(interval) * time.Second,
			ChatDescriptionTemplate: chatDescriptionTemplate,
		},
	)
	bot.Run()
}
