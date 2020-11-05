package bot

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gowee/github-status-bot/pkg/api"
	"github.com/gowee/github-status-bot/pkg/data"
	"github.com/gowee/github-status-bot/pkg/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

var statusPage = api.StatusPage{
	PageID: api.GitHubPageID,
} // TODO: move into Bot

func (bot *Bot) trackUpdates(stop chan struct{}) {
	log.Println("Check interval:", bot.CheckInterval)
	bot.updateOnce(true)
	tick := time.Tick(bot.CheckInterval)
	for {
		select {
		case <-tick:
			// log.Println("tick!")
			bot.updateOnce(false)
			break
		case <-stop:
			// close(stop)
			log.Println("Track task stops")
			return
		}
	}
}

func (bot *Bot) updateOnce(forceUpdate bool) {
	// if forceUpdate {
	// 	log.Println("Force updating (usually the first interval)")
	// }
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

	// Do not force update chat photo/description as its service message cannot be silenced
	if currSts.Status.Indicator != prevData.GlobalStatusIndicator {
		err := bot.Client.SetGroupPhoto(bot.Chat,
			&tb.Photo{File: tb.File{FileReader: currSts.Status.ToIcon()}})
		if err == nil {
			log.Println("Updated chat photo")
		} else {
			log.Println("Failed to update chat photo: ", err)
		}
		err = bot.Client.SetGroupTitle(
			bot.Chat,
			fmt.Sprintf("%sGitHub: %s", currSts.Status.ToEmoji(), currSts.Status.Description),
		)
		if err == nil {
			log.Println("Updated chat title")
		} else {
			log.Println("Failed to update chat title: ", err)
		}
	}

	if bot.ChatDescriptionTemplate != "" &&
		(forceUpdate || currSts.Status.Indicator != prevData.GlobalStatusIndicator) {
		err := bot.Client.SetGroupDescription(
			bot.Chat,
			fmt.Sprintf(bot.ChatDescriptionTemplate, formatMultipleComponents(currSts.Components)),
		)
		if err == nil {
			log.Println("Updated chat description")
		} else if !strings.Contains(err.Error(), ": chat description is not modified") {
			log.Println("Failed to update chat description: ", err)
		}
		prevData.GlobalStatusIndicator = currSts.Status.Indicator
	}

	// newEvents := []api.Event{}
	var newEvents []api.Event

	// WTF: why no generic Min/Max? even if compiler tricks would be great!
	for idx := range currSts.Incidents[0:utils.Min(10, len(currSts.Incidents))] {
		// WTF: why for...range... cannot get references?
		//   ref: https://stackoverflow.com/questions/20185511/range-references-instead-values
		newEvents = append(newEvents, &currSts.Incidents[idx])
	}
	for idx := range currSts.ScheduledMaintenances[0:utils.Min(5, len(currSts.ScheduledMaintenances))] {
		newEvents = append(newEvents, &currSts.ScheduledMaintenances[idx])
	}
	sort.Sort(api.Events(newEvents))

	for _, event := range newEvents {
		if prevEvent, ok := prevData.Events[event.GetID()]; ok {
			bot.updateEvent(prevEvent, event, forceUpdate)
		} else {
			storedEvent := bot.newEvent(event)
			if event != nil {
				prevData.Events[event.GetID()] = storedEvent
			}
		}
	}

	// TODO: transactionally update status
	if err := bot.DB.Store(prevData); err != nil {
		log.Fatal("Failed to update database: ", err)
		return
	}
}

func (bot *Bot) newEvent(newEvent api.Event) *data.Event {
	silent := !newEvent.ShouldNotify()
	msg, err := bot.Client.Send(
		tb.ChatID(bot.Chat.ID),
		newEvent.Format(),
		&tb.SendOptions{
			DisableWebPagePreview: true,
			DisableNotification:   silent,
			ParseMode:             "HTML",
		},
	)
	if err != nil {
		log.Println("Failed to send message for event: ", newEvent.GetID(), err)
		return nil
	}
	return &data.Event{
		ID:               newEvent.GetID(),
		UpdatedAt:        newEvent.GetUpdatedAt(),
		MessageReference: tb.StoredMessage{MessageID: strconv.Itoa(msg.ID), ChatID: bot.Chat.ID},
	}
}

func (bot *Bot) updateEvent(storedEvent *data.Event, newEvent api.Event, forceUpdate bool) {
	silent := !newEvent.ShouldNotify() // Will it work for editMessage?
	if forceUpdate || newEvent.GetUpdatedAt().After(storedEvent.UpdatedAt) {
		_, err := bot.Client.Edit(
			storedEvent.MessageReference,
			newEvent.Format(),
			&tb.SendOptions{
				DisableWebPagePreview: true,
				DisableNotification:   silent,
				ParseMode:             "HTML",
			},
		)
		// WTF: why telebot's error is not typed so that to compared with errors.Is?

		// err != tb.ErrMessageNotModified does not work due to a bug.
		// 	ref: https://github.com/tucnak/telebot/issues/330

		// WTF: why no string.Contains(needle)?
		if err == nil {
			log.Println("Updated incident: ", newEvent.GetID())
		} else if !strings.Contains(err.Error(), ": message is not modified") {
			log.Println("Failed to update message for incident: ", newEvent.GetID(), err)
			return
		} // else: No change.
		storedEvent.UpdatedAt = newEvent.GetUpdatedAt()
	} // else not new
}
