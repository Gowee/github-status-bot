module github.com/gowee/github-status-bot

go 1.15

require (
	github.com/gobuffalo/packr v1.30.1
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/tools v0.0.0-20201103235415-b653051172e4 // indirect
	golang.org/x/tools/gopls v0.5.2 // indirect
	gopkg.in/tucnak/telebot.v2 v2.3.5
)

// Note: clean this via `go mod tidy`
// WTF: why `go get` a unrelated binary package resulting in an entry here?
