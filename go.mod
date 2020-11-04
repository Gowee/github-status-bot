module github.com/gowee/github-status-bot

go 1.15

require (
	github.com/gobuffalo/packr v1.30.1
	github.com/rogpeppe/go-internal v1.6.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/tucnak/telebot.v2 v2.3.5
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

// Note: clean this via `go mod tidy`
// WTF: why `go get` a unrelated binary package resulting in an entry here?
