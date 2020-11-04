package data

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var EmtpyRootDocument = RootDocument{
	Events:     make(map[string]*Event),
	Components: make(map[string]*Component),
}

type RootDocument struct {
	GlobalStatusIndicator string                `json:"global_status_indicator"`
	Events                map[string]*Event     `json:"events"`
	Components            map[string]*Component `json:"component"`
}

// An Event is either a incident or a scheduled maintenance, associated with a Telegram message.
type Event struct {
	ID               string           `json:"id"`
	UpdatedAt        time.Time        `json:"updated"`
	MessageReference tb.StoredMessage `json:"message_reference"`
}

type Component struct {
	// Not yet used
	ID        string    `json:"id"`
	Indicator string    `json:"indicator"`
	UpdatedAt time.Time `json:"updated_at"`
}
