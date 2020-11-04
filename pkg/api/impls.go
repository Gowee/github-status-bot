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
	statusIcon := i.ToImpactEmoji() // "⚠️"
	if i.Status == "resolved" {
		statusIcon = "✅"
	}
	header := fmt.Sprintf("<b><a href=\"%s\">%s</a></b> %s\n", i.Shortlink, i.Name, statusIcon)

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

func (i *Incident) ShouldNotify() bool {
	switch i.Impact {
	case "none":
		fallthrough
	case "minor":
		return false
	case "major":
		fallthrough
	case "critical":
		return true
	default:
		log.Println("Unknown incident impact: ", i.Impact)
		return true
	}
}

func (i *Incident) ToImpactEmoji() string {
	impact := i.Impact
	switch impact {
	case "none":
		return ""
	case "minor":
		return "❕"
	case "major":
		return "❗️"
	case "critical":
		return "‼️"
	default:
		log.Println("Unknown incident impact: ", impact)
		return "❔"
	}
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
		return ""
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

func (c *Component) ToStatusEmoji() string {
	switch c.Status {
	case "operational":
		return "✅"
	case "degraded_performance":
		return "❕"
	case "partial_outage":
		return "❗️"
	case "major_outage":
		return "‼️"
	default:
		log.Printf("Unknown status: %s, for component: %s\n", c.Status, c.Name)
		return "❔"
	}
}

func (c *Component) ToStatusSimple() string {
	// Ref: https://www.githubstatus.com/ source#L1532
	switch c.Status {
	case "operational":
		return "Normal"
	case "degraded_performance":
		return "Degraded"
	case "partial_outage":
		return "Degraded"
	case "major_outage":
		return "Incident"
	default:
		log.Printf("Unknown status: %s, for component: %s\n", c.Status, c.Name)
		return "Unknown"
	}
}

func (c *Component) Format() string {
	return fmt.Sprintf("%s: %s%s", c.Name, c.ToStatusSimple(), c.ToStatusEmoji())
}

// WTF: why go fmt does not break long lines (by default?)?