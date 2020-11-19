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
	chatServiceMessages     chan time.Time // service messages to be deleted
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
		panic(
			fmt.Sprintf(
				"chatDescriptionTemplate is present but invalid: %s",
				chatDescriptionTemplate,
			),
		)
	}

	return Bot{
		Client:                  client,
		Chat:                    chat,
		DB:                      &db,
		CheckInterval:           interval,
		ChatDescriptionTemplate: options.ChatDescriptionTemplate,
		chatServiceMessages:     make(chan time.Time),
	}
}

func (bot *Bot) Run() {
	log.Println("Up and running...")
	log.Println("  for:", renderChat(bot.Chat))
	log.Println("  as: ", renderUser(bot.Client.Me))
	bot.Client.Handle("/hello", func(m *tb.Message) {
		bot.Client.Send(m.Sender, "Hello World!")
	})

	// delete service messages of chat title/description updates
	bot.Client.Handle(tb.OnChannelPost, func(m *tb.Message) {
		if m.Chat.ID == bot.Chat.ID {
			if m.NewGroupPhoto != nil || m.NewGroupTitle != "" {
				// relate setChatTitle/Photo action with service message by time
				for {
					timeout := time.After(3 * time.Second)
					select {
					case t := <-bot.chatServiceMessages:
						if time.Now().Sub(t) > 3*time.Second {
							// a possible case:
							// setChatTitle returns True, resulting a entry in the channel, while
							// the new title is the same as the previous one, resulting in no
							// (new) service message
							log.Println("Discarding a stale entry in chatServiceMessages channel")
						} else {
							if err := bot.Client.Delete(m); err == nil {
								log.Println("Deleted a service message")
							} else {
								log.Println("Failed to delete a service message", err)
							}
							return
						}
						break
					case <-timeout:
						log.Println(
							"A service message is not deleted as it seems to be sent by others",
						)
						return
						break
					}
				}
			}
		}
	})

	stop := make(chan struct{})
	go bot.trackUpdates(stop)

	bot.Client.Start()

	close(stop)
}
