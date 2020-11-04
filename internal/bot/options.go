package bot

import (
	"github.com/urfave/cli/v2" 
)

var Options struct {
	BotToken      string `arg:"positional,required,env:BOT_TOKEN"`
	ChatID        int64  `arg:"positional,required,env:CHAT_ID" help:"the numeric ID of a chat to send updates"`
	DataFile      string `arg:"-d,--data-file" default:"./data.json" help:"path to the data file"`
	CheckInterval int    `arg:"-t,--check-interval,env:CHECK_INTERVAL" default:"300" help:"check interval in seconds"`
}

func GetOptions() Options {
	
}