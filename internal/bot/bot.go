package bot

import (
	"log"
	"time"

	"github.com/gowee/github-status-bot/pkg/data"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Run(options Options) {
	if err := data.EnsureInitialized(); err != nil {
		log.Fatal(err)
		return
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  options.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Up and running...")

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello World!")
	})

	monitorStop := make(chan struct{})
	go trackUpdates(b, options.ChatID, monitorStop)

	b.Start()

	close(monitorStop)
}
