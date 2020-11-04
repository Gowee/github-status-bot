package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gowee/github-status-bot/pkg/api"
	"github.com/gowee/github-status-bot/pkg/data"
	"github.com/gowee/github-status-bot/pkg/utils"
	tb "gopkg.in/tucnak/telebot.v2"
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

	if err := data.EnsureInitialized(); err != nil {
		log.Fatal(err)
		return
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
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
	go trackUpdates(b, chatID, monitorStop)

	b.Start()

	close(monitorStop)
}

func trackUpdates(bot *tb.Bot, chatID int64, stop chan struct{}) {
	statusPage := api.StatusPage{
		PageID: "kctbh9vrtdwd",
	}
	tick := time.Tick(5 * time.Second)
	for {
		select {
		case <-tick:
			// fmt.Println("tick!")
			allSts, err := statusPage.QueryOverall()
			if err != nil {
				log.Println(err)
				continue
			}
			prevData, err := data.Load()
			if err != nil {
				log.Fatal("Failed to read database: ", err)
				return
			}

			if allSts.Status.Indicator != prevData.GlobalStatusIndicator {
				err := bot.SetGroupPhoto(&tb.Chat{ID: chatID},
					&tb.Photo{File: tb.File{FileReader: allSts.Status.ToIcon()}})
				if err != nil {
					log.Println("Failed to update chat photo: ", err)
				}
				err = bot.SetGroupTitle(&tb.Chat{ID: chatID}, fmt.Sprintf("GitHub: %s %s", allSts.Status.Description, allSts.Status.ToEmoji()))
				if err != nil {
					log.Println("Failed to update chat title: ", err)
				}

				prevData.GlobalStatusIndicator = allSts.Status.Indicator
			}

			// WTF: why no generic Min/Max? even if compiler tricks would be great!
			for _, incident := range allSts.Incidents[0:utils.Min(10, len(allSts.Incidents))] {
				if prevEvent, ok := prevData.Events[incident.ID]; ok {
					if incident.UpdatedAt.After(prevEvent.UpdatedAt) {
						_, err := bot.Edit(prevEvent.MessageReference, incident.Format(), &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "HTML"})
						// WTF: why telebot's error is not typed so that to compared with errors.Is?

						// err != tb.ErrMessageNotModified does not work due to a bug.
						// 	ref: https://github.com/tucnak/telebot/issues/330

						// WTF: why no string.Contains(needle)?
						if err == nil {
							log.Println("Updated incident: ", incident.ID)
						} else if !strings.Contains(err.Error(), ": message is not modified:") {
							log.Println("Failed to update message for incident: ", incident.ID, err)
							continue
						} // else: No change.
						prevEvent.UpdatedAt = incident.UpdatedAt
					}
				} else {
					msg, err := bot.Send(tb.ChatID(chatID), incident.Format(), &tb.SendOptions{DisableWebPagePreview: true, ParseMode: "HTML"})
					if err != nil {
						log.Println("Failed to send message for incident: ", incident.ID, err)
						continue
					}
					prevData.Events[incident.ID] = &data.Event{
						ID:               incident.ID,
						UpdatedAt:        incident.UpdatedAt,
						MessageReference: tb.StoredMessage{MessageID: strconv.Itoa(msg.ID), ChatID: chatID},
					}
					log.Println("New incident: ", incident.ID)
				}
			}
			if err := data.Store(prevData); err != nil {
				log.Fatal("Failed to update database: ", err)
				return
			}
			break
		case <-stop:
			// close(stop)
			log.Println("Monitor task stops.")
			return
		}
	}
}
