package bot

import (
	"fmt"
	"strings"

	"github.com/gowee/github-status-bot/pkg/api"
	tb "gopkg.in/tucnak/telebot.v2"
)

func formatMultipleComponents(components []api.Component) string {
	seq := make([]string, 0)
	for _, component := range components {
		if component.ID != api.IgnoredGitHubDummyComponentID {
			// skip one entry
			seq = append(seq, component.Format())
		}
	}
	return strings.Join(seq, "\n")
}

func renderChat(chat *tb.Chat) string {
	// Ref: https://core.telegram.org/bots/api#chat
	text := ""
	if chat.Type == "private" {
		if chat.FirstName != "" && chat.LastName != "" {
			text = fmt.Sprintf("%s %s", chat.FirstName, chat.LastName)
		} else if chat.FirstName != "" {
			text = chat.FirstName
		} else if chat.LastName != "" {
			text = chat.LastName
		} else {
			text = fmt.Sprint(chat.ID)
		}
	} else {
		text = chat.Title
	}
	if chat.Username != "" {
		text += fmt.Sprintf(" (@%s)", chat.Username)
	}
	return text
}

func renderUser(user *tb.User) string {
	// Ref: https://core.telegram.org/bots/api#user
	text := ""
	if user.FirstName != "" && user.LastName != "" {
		text = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	} else if user.FirstName != "" {
		text = user.FirstName
	} else if user.LastName != "" {
		text = user.LastName
	} else {
		text = fmt.Sprint(user.ID)
	}
	if user.Username != "" {
		text += fmt.Sprintf(" (@%s)", user.Username)
	}
	return text
}
