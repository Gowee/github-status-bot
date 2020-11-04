package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/gowee/github-status-bot/pkg/assets"
	"github.com/gowee/github-status-bot/pkg/utils"
)

const updateDateLayout = "2006-01-02 15:04:05"

func (i *Incident) Format() string {
	statusIcon := "⚠️"
	if i.Status == "resolved" {
		statusIcon = "✅"
	}
	header := fmt.Sprintf("%s <b><a href=\"%s\">%s</a></b>\n", statusIcon, i.Shortlink, i.Name)

	lines := make([]string, 1+utils.Min(3+1, len(i.IncidentUpdates)))
	lines[0] = header
	// WTF: why no dynamic array?
	// 	ref: https://stackoverflow.com/questions/33834742/remove-and-adding-elements-to-array-in-go-lang
	// 	ref: https://ewencp.org/blog/golang-iterators/index.html
	// WTF: why no combinator such as map?
	if len(i.IncidentUpdates) > 3 {
		// The original updates are sorted descendingly by date.
		lines[1] = i.IncidentUpdates[len(i.IncidentUpdates)-1].Format()
		lines[2] = fmt.Sprintf("<pre>----- %d update omitted -----</pre>", len(i.IncidentUpdates)-3)
		lines[3] = i.IncidentUpdates[1].Format()
		lines[4] = i.IncidentUpdates[0].Format()
	} else {
		for idx, update := range i.IncidentUpdates {
			lines[len(i.IncidentUpdates)-idx] = update.Format()
		}
		// WTF: why no built-in reverse?

	}
	return strings.Join(lines, "\n")
}

func (u *IncidentUpdate) Format() string {
	// WTF: why no format literal?
	return fmt.Sprintf("<b>%s</b> <i>at %s</i>:\n<u>%s</u>",
		u.Status,
		u.UpdatedAt.Format(updateDateLayout),
		u.Body)
}

func (s *Status) ToIcon() io.Reader {
	indicator := s.Indicator
	switch indicator {
	case "none":
		return bytes.NewReader(assets.GitHubIconNormal)
	case "minor":
		return bytes.NewReader(assets.GitHubIconYellow)
	case "major":
		return bytes.NewReader(assets.GitHubIconOrange)
	case "critical":
		return bytes.NewReader(assets.GitHubIconYellow)
	default:
		log.Println("Unknown status indicator: ", indicator)
		return bytes.NewReader(assets.GitHubIconNormal)
	}
}

func (s *Status) ToEmoji() string {
	indicator := s.Indicator
	switch indicator {
	case "none":
		return "✅"
	case "minor":
		return "❕"
	case "major":
		return "❗️"
	case "critical":
		return "‼️"
	default:
		log.Println("Unknown status indicator: ", indicator)
		return "❔"
	}
}

// WTF: why go fmt does not break long lines (by default?)?
