package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gowee/github-status-bot/pkg/data"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	Client                  *tb.Bot
	Chat                    *tb.Chat
	DB                      *data.Database
	CheckInterval           time.Duration
	ChatDescriptionTemplate string
}

func NewBotFromOptions(options Options) Bot {
	db := data.NewDBFromFilePath(options.DataFilePath)
	if err := db.EnsureInitialized(); err != nil {
		log.Fatal(err)
	}

	client, err := tb.NewBot(tb.Settings{
		Token:  options.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal("Failed to create a telebot Bot", err)
	}

	chat, err := client.ChatByID(options.ChatID)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to the chat %s with error %s", options.ChatID, err))
	}

	interval := options.CheckInterval
	if interval <= 0*time.Second {
		interval = 5 * time.Minute // default interval
	} else if interval < 5*time.Second {
		interval = 5 * time.Second // minimum interval
	}
	// WTF: why the auto-enforced code format for * is different here?

	chatDescriptionTemplate := options.ChatDescriptionTemplate
	if chatDescriptionTemplate != "" && !strings.Contains(chatDescriptionTemplate, "%s") {
		fmt.Println(chatDescriptionTemplate)
		panic("chatDescriptionTemplate is present but invalid")
	}

	return Bot{
		Client:                  client,
		Chat:                    chat,
		DB:                      &db,
		CheckInterval:           interval,
		ChatDescriptionTemplate: options.ChatDescriptionTemplate,
	}
}

func (bot *Bot) Run() {
	log.Println("Up and running...")
	log.Println("  for:", renderChat(bot.Chat))
	log.Println("  as: ", renderUser(bot.Client.Me))
	bot.Client.Handle("/hello", func(m *tb.Message) {
		bot.Client.Send(m.Sender, "Hello World!")
	})

	// bot.Client.Handle(tb.OnChannelPost, func(m *tb.Message) {
	// 	if m.Chat.ID == bot.Chat.ID {
	// 		// channel posts only
	// 		// log.Println("channel post", m)
	// 		if m.NewGroupPhoto != nil {
	// 			if err := bot.Client.Delete(m); err == nil {
	// 				log.Println("Deleted a NewGroupPhoto message")
	// 			} else {
	// 				log.Println("Failed to delete a NewGroupPhoto message", err)
	// 			}
	// 		} else if m.NewGroupTitle != "" {
	// 			if err := bot.Client.Delete(m); err == nil {
	// 				log.Println("Deleted a NewGroupTitle message")
	// 			} else {
	// 				log.Println("Failed to delete a NewGroupTitle message", err)
	// 			}
	// 		}
	// 	}
	// })

	stop := make(chan struct{})
	go bot.trackUpdates(stop)

	bot.Client.Start()

	close(stop)
}
