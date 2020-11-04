package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gowee/github-status-bot/pkg/api"
	"github.com/gowee/github-status-bot/pkg/data"
	"github.com/gowee/github-status-bot/pkg/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

var statusPage = api.StatusPage{
	PageID: "kctbh9vrtdwd",
} // TODO: move into Bot

func (bot *Bot) updateOnce() {
	currSts, err := statusPage.QueryOverall()
	if err != nil {
		log.Println(err)
		return
	}
	prevData, err := bot.DB.Load()
	if err != nil {
		log.Fatal("Failed to read database: ", err)
		return
	}

	if currSts.Status.Indicator != prevData.GlobalStatusIndicator {
		err := bot.Client.SetGroupPhoto(bot.Chat,
			&tb.Photo{File: tb.File{FileReader: currSts.Status.ToIcon()}})
		if err != nil {
			log.Println("Failed to update chat photo: ", err)
		}
		err = bot.Client.SetGroupTitle(
			bot.Chat,
			fmt.Sprintf("GitHub: %s %s", currSts.Status.Description, currSts.Status.ToEmoji()),
		)
		if err != nil {
			log.Println("Failed to update chat title: ", err)
		}

		prevData.GlobalStatusIndicator = currSts.Status.Indicator
	}

	// WTF: why no generic Min/Max? even if compiler tricks would be great!
	for _, incident := range currSts.Incidents[0:utils.Min(10, len(currSts.Incidents))] {
		if prevEvent, ok := prevData.Events[incident.ID]; ok {
			if incident.UpdatedAt.After(prevEvent.UpdatedAt) {
				_, err := bot.Client.Edit(
					prevEvent.MessageReference,
					incident.Format(),
					&tb.SendOptions{DisableWebPagePreview: true, ParseMode: "HTML"},
				)
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
			msg, err := bot.Client.Send(
				tb.ChatID(bot.Chat.ID),
				incident.Format(),
				&tb.SendOptions{DisableWebPagePreview: true, ParseMode: "HTML"})
			if err != nil {
				log.Println("Failed to send message for incident: ", incident.ID, err)
				continue
			}
			prevData.Events[incident.ID] = &data.Event{
				ID:               incident.ID,
				UpdatedAt:        incident.UpdatedAt,
				MessageReference: tb.StoredMessage{MessageID: strconv.Itoa(msg.ID), ChatID: bot.Chat.ID},
			}
			log.Println("New incident: ", incident.ID)
		}
	}
	if err := bot.DB.Store(prevData); err != nil {
		log.Fatal("Failed to update database: ", err)
		return
	}
}

func (bot *Bot) trackUpdates(stop chan struct{}) {
	log.Println("Track task starts with interval", bot.CheckInterval)
	tick := time.Tick(bot.CheckInterval)
	for {
		select {
		case <-tick:
			log.Println("tick!")
			bot.updateOnce()
			break
		case <-stop:
			// close(stop)
			log.Println("Track task stops")
			return
		}
	}
}
