module github.com/gowee/github-status-bot

go 1.15

require (
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/rogpeppe/go-internal v1.6.2 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/cobra v1.1.1 // indirect
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9 // indirect
	golang.org/x/sys v0.0.0-20201101102859-da207088b7d1 // indirect
	golang.org/x/tools v0.0.0-20201103235415-b653051172e4 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/tucnak/telebot.v2 v2.3.5
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

// Note: clean this via `go mod tidy`
// WTF: why `go get` a unrelated binary package resulting in an entry here?
