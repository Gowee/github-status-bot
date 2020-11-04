package bot

import (
	"time"
)

// Options for bot, usually extracted from cli args or env vars.
type Options struct {
	BotToken      string
	ChatID        int64
	DataFilePath  string
	CheckInterval time.Duration
}

// TODO: extract options with arg parser?
