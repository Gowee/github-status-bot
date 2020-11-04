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
	log.Println("Track task starts with interval", bot.CheckInterval)
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
	if forceUpdate {
		log.Println("Force updating (usually the first interval)")
	}
	currSts, err := statusPage.QueryOverall()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Events fetched")
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

	if forceUpdate || currSts.Status.Indicator != prevData.GlobalStatusIndicator {
		err = bot.Client.SetGroupDescription(bot.Chat, formatMultipleComponents(currSts.Components))
		if err == nil {
			log.Println("Updated chat description")
		} else if !strings.Contains(err.Error(), ": chat description is not modified") {
			log.Println("Failed to update chat description: ", err)
		}
		prevData.GlobalStatusIndicator = currSts.Status.Indicator
	}
	log.Println("Global updated")

	// newEvents := []api.Event{}
	var newEvents []api.Event

	// WTF: why no generic Min/Max? even if compiler tricks would be great!
	for idx, _ := range currSts.Incidents[0:utils.Min(10, len(currSts.Incidents))] {
		// WTF: where for...range... cannot get references?
		//   ref: https://stackoverflow.com/questions/20185511/range-references-instead-values
		newEvents = append(newEvents, &currSts.Incidents[idx])
		// if prevEvent, ok := prevData.Events[incident.ID]; ok {
		// 	bot.updateEvent(prevEvent, &incident, forceUpdate)
		// } else {
		// 	event := bot.newEvent(&incident)
		// 	if event != nil {
		// 		prevData.Events[incident.ID] = event
		// 	}
		// }
	}
	for idx, _ := range currSts.ScheduledMaintenances[0:utils.Min(5, len(currSts.ScheduledMaintenances))] {
		newEvents = append(newEvents, &currSts.ScheduledMaintenances[idx])
		// if prevEvent, ok := prevData.Events[maintenance.ID]; ok {
		// 	bot.updateEvent(prevEvent, &maintenance, forceUpdate)
		// } else {
		// 	event := bot.newEvent(&maintenance)
		// 	if event != nil {
		// 		prevData.Events[maintenance.ID] = event
		// 	}
		// }
	}
	log.Println(newEvents)
	sort.Sort(api.Events(newEvents))
	log.Println("Events sorted", len(newEvents))

	for _, event := range newEvents {
		log.Println("iterating at ", event.GetID())
		if prevEvent, ok := prevData.Events[event.GetID()]; ok {
			log.Println("branch 1")
			bot.updateEvent(prevEvent, event, forceUpdate)
		} else {
			log.Println("branch 2")
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
	log.Println("New event start: ", newEvent.GetID())
	silent := !newEvent.ShouldNotify()
	msg, err := bot.Client.Send(
		tb.ChatID(bot.Chat.ID),
		newEvent.Format(),
		&tb.SendOptions{DisableWebPagePreview: true, DisableNotification: silent, ParseMode: "HTML"})
	log.Println("New event sent: ", newEvent.GetID())
	if err != nil {
		log.Println("Failed to send message for event: ", newEvent.GetID(), err)
		return nil
	}
	log.Println("New event: ", newEvent.GetID())
	return &data.Event{
		ID:               newEvent.GetID(),
		UpdatedAt:        newEvent.GetUpdatedAt(),
		MessageReference: tb.StoredMessage{MessageID: strconv.Itoa(msg.ID), ChatID: bot.Chat.ID},
	}
}

func (bot *Bot) updateEvent(storedEvent *data.Event, newEvent api.Event, forceUpdate bool) {
	log.Println("Update event start: ", newEvent.GetID())
	silent := !newEvent.ShouldNotify() // Will it work for editMessage?
	if forceUpdate || newEvent.GetUpdatedAt().After(storedEvent.UpdatedAt) {
		_, err := bot.Client.Edit(
			storedEvent.MessageReference,
			newEvent.Format(),
			&tb.SendOptions{DisableWebPagePreview: true, DisableNotification: silent, ParseMode: "HTML"},
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
