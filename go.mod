module github.com/gowee/github-status-bot

go 1.15

require (
	github.com/gobuffalo/packr v1.30.1
	gopkg.in/tucnak/telebot.v2 v2.3.5
)

// Note: clean this via `go mod tidy`
// WTF: why `go get` a unrelated binary package resulting in an entry here?
 