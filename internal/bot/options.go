package bot

import (
	"time"
)

// Options for bot, usually extracted from cli args or env vars.
type Options struct {
	BotToken                string
	ChatID                  string
	DataFilePath            string
	CheckInterval           time.Duration
	ChatDescriptionTemplate string
}

// TODO: extract options with arg parser?
