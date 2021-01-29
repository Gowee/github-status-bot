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

const updateDateLayout = "Jan 2, 15:04" // 2006-01-02 15:04:05
const maximumUpdatesPerMessage = 0xFF   // MUST >= 2

func formatIncidentOrScheduledMaintenance(name string, url string, statusIcon string, updates []IncidentUpdate) string {
	header := fmt.Sprintf("<b>%s</b> <a href=\"%s\">%s</a>\n\n", name, url, statusIcon)
	// if i.Status != "resolved" {
	// 	header += fmt.Sprintf("(%s)", i.Status)
	// }

	lines := make([]string, utils.Min(maximumUpdatesPerMessage+1, len(updates)))
	// W-T-F: why no dynamic array?
	// 	ref: https://stackoverflow.com/questions/33834742/remove-and-adding-elements-to-array-in-go-lang
	// 	ref: https://ewencp.org/blog/golang-iterators/index.html
	// WTFUpdate: It is: https://tour.golang.org/moretypes/13
	// 	 note: the append pattern just exposes the release/alloc mem prodecure as is in other languages
	// WTF: why no combinator such as map?

	// The original updates are sorted descendingly by date. Just reverse it.
	if len(updates) > maximumUpdatesPerMessage {
		lines[0] = updates[len(updates)-1].Format("")
		suf := ""
		if len(updates) > maximumUpdatesPerMessage+1 {
			suf = "s"
		}
		lines[1] = fmt.Sprintf(
			"<pre>┄┄┄┄┄ %d update%s omitted ┄┄┄┄┄</pre>",
			len(updates)-maximumUpdatesPerMessage,
			suf,
		)
		for i, update := range updates[1 : maximumUpdatesPerMessage-1] {
			lines[len(lines)-1-i-1] = update.Format("")
		}
		lines[len(lines)-1] = updates[0].Format(url)
	} else {
		lines[len(updates)-1] = updates[0].Format(url)
		for i, update := range updates[1:] {
			lines[len(updates)-1-i-1] = update.Format("")
		}
		// WTF: why no built-in reverse?
		// WTFUpdate: there is, but is hard to use due to the poor type system
		//	 ref: https://stackoverflow.com/a/18343326/5488616
	}
	return header + strings.Join(lines, "\n\n")
}

func (i *Incident) Format() string {
	statusIcon := i.ToImpactEmoji() // "⚠️"
	// If the status indicates that the indident has finished, then the impact does not matter.
	if i.Status == "resolved" {
		statusIcon = "✅"
	} else if i.Status == "postmortem" {
		statusIcon = "☑️"
	}
	return formatIncidentOrScheduledMaintenance(i.Name, i.Shortlink, statusIcon, i.IncidentUpdates)
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

func (i *Incident) IsFinished() bool {
	switch i.Status {
	case "resolved":
		fallthrough
	case "postmortem":
		return true
	default:
		return false
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

func (sm *ScheduledMaintenance) Format() string {
	statusIcon := sm.ToStatusEmoji()
	// Currently the impact is not showed.
	return formatIncidentOrScheduledMaintenance(sm.Name, sm.Shortlink, statusIcon, sm.IncidentUpdates)
}

func (sm *ScheduledMaintenance) ToStatusEmoji() string {
	status := sm.Status
	switch status {
	case "scheduled":
		return "ℹ️"
	case "in_progress":
		return "⏳"
	case "verifying":
		return "⌛"
	case "completed":
		// Some maintenance event indicates a break change which might not be fine.
		// So here use a check mark different to resolved incidents to distinguish.
		return "☑️"
	default:
		log.Println("Unknown scheduled maintenance status: ", status)
		return "❔"
	}
}

func (sm *ScheduledMaintenance) ShouldNotify() bool {
	switch sm.Impact {
	case "none":
		fallthrough
	case "minor":
		return false
	case "major":
		fallthrough
	case "critical":
		fallthrough
	case "maintenance":
		return true
	default:
		log.Println("Unknown ScheduledMaintenance impact: ", sm.Impact)
		return true
	}
}

func (sm *ScheduledMaintenance) IsFinished() bool {
	switch sm.Status {
	case "completed":
		return true
	default:
		return false
	}
}

func (u *IncidentUpdate) Format(url string) string {
	// WTF: why no format literal?
	var status string = u.Status
	if url != "" {
		status = fmt.Sprintf("<a href=\"%s\">%s</a>", url, u.Status)
	}
	return fmt.Sprintf("<b>%s</b> - %s <i>@ %s</i>",
		status,
		u.Body,
		u.UpdatedAt.Format(updateDateLayout))
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
	case "maintenance":
		return bytes.NewReader(assets.GitHubIconBlue) // Undocumented status indicator
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
	case "maintenance":
		return "🛠️" // Undocumented status indicator
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
	case "under_maintenance": // Undocumented component status
		return "🛠️"
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
	case "under_maintenance": // Undocumented component status
		return "Maintenance"
	default:
		log.Printf("Unknown status: %s, for component: %s\n", c.Status, c.Name)
		return "Unknown"
	}
}

func (c *Component) Format() string {
	return fmt.Sprintf("%s: %s %s", c.Name, c.ToStatusSimple(), c.ToStatusEmoji())
}

// WTF: why go fmt does not break long lines (by default?)?
// WTFUpdate: there is a 3rd project called golines, but it seems not hot
