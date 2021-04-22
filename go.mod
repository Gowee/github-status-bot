module github.com/gowee/github-status-bot

go 1.15

require (
	github.com/gobuffalo/packr/v2 v2.8.1
	github.com/gomarkdown/markdown v0.0.0-20210408062403-ad838ccf8cdd
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210421221651-33663a62ff08 // indirect
	golang.org/x/term v0.0.0-20210421210424-b80969c67360 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/tucnak/telebot.v2 v2.3.5
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// Note: clean this via `go mod tidy`
// WTF: why `go get` a unrelated binary package resulting in an entry here?
// WTF: why no dev dependencies?
