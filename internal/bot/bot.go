package bot

import (
	"log"
	"time"

	"github.com/gowee/github-status-bot/pkg/data"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	Client        *tb.Bot
	Chat          *tb.Chat
	DB            *data.Database
	CheckInterval time.Duration
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
		log.Fatal(err)
	}

	interval := options.CheckInterval
	if interval <= 0*time.Second {
		interval = 5 * time.Minute // default interval
	} else if interval < 5*time.Second {
		interval = 5 * time.Second // minimum interval
	}
	// WTF: why the auto-enforced format for * is different here?

	return Bot{
		Client:        client,
		Chat:          &tb.Chat{ID: options.ChatID},
		DB:            &db,
		CheckInterval: interval,
	}
}

func (bot *Bot) Run() {
	log.Println("Up and running...")

	bot.Client.Handle("/hello", func(m *tb.Message) {
		bot.Client.Send(m.Sender, "Hello World!")
	})

	stop := make(chan struct{})
	go bot.trackUpdates(stop)

	bot.Client.Start()

	close(stop)
}
