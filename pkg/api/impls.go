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

func (i *Incident) Format() string {
	statusIcon := i.ToImpactEmoji() // "⚠️"
	// If the status indicates that the indident has finished, then the impact does not matter.
	if i.Status == "resolved" {
		statusIcon = "✅"
	} else if i.Status == "postmortem" {
		statusIcon = "☑️"
	}
	header := fmt.Sprintf("<b>%s</b> <a href=\"%s\">%s</a>\n", i.Name, i.Shortlink, statusIcon)
	// if i.Status != "resolved" {
	// 	header += fmt.Sprintf("(%s)", i.Status)
	// }

	lines := make([]string, 1+utils.Min(3+1, len(i.IncidentUpdates)))
	lines[0] = header
	// W-T-F: why no dynamic array?
	// 	ref: https://stackoverflow.com/questions/33834742/remove-and-adding-elements-to-array-in-go-lang
	// 	ref: https://ewencp.org/blog/golang-iterators/index.html
	// WTFUpdate: It is: https://tour.golang.org/moretypes/13
	// 	 note: the append pattern just exposes the release/alloc mem prodecure as is in other languages
	// WTF: why no combinator such as map?
	if len(i.IncidentUpdates) > 3 {
		// The original updates are sorted descendingly by date.
		lines[1] = i.IncidentUpdates[len(i.IncidentUpdates)-1].Format("")
		suf := ""
		if len(i.IncidentUpdates) > 4 {
			suf = "s"
		}
		lines[2] = fmt.Sprintf("<pre>┄┄┄┄┄ %d update%s omitted ┄┄┄┄┄</pre>", len(i.IncidentUpdates)-3, suf)
		lines[3] = i.IncidentUpdates[1].Format("")
		lines[4] = i.IncidentUpdates[0].Format(i.Shortlink)
	} else {
		for idx, update := range i.IncidentUpdates {
			url := ""
			if idx == 0 {
				url = i.Shortlink
			}
			lines[len(i.IncidentUpdates)-idx] = update.Format(url)
		}
		// WTF: why no built-in reverse?
		// WTFUpdate: there is, but is hard to use due to the poor type system
		//	 ref: https://stackoverflow.com/a/18343326/5488616
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
	header := fmt.Sprintf("<b>%s</b> <a href=\"%s\">%s</a>\n", sm.Name, sm.Shortlink, statusIcon)

	lines := make([]string, 1+utils.Min(3+1, len(sm.IncidentUpdates)))
	lines[0] = header
	if len(sm.IncidentUpdates) > 3 {
		lines[1] = sm.IncidentUpdates[len(sm.IncidentUpdates)-1].Format("")
		suf := ""
		if len(sm.IncidentUpdates) > 4 {
			suf = "s"
		}
		lines[2] = fmt.Sprintf("<pre>┄┄┄┄┄ %d update%s omitted ┄┄┄┄┄</pre>", len(sm.IncidentUpdates)-3, suf)
		lines[3] = sm.IncidentUpdates[1].Format("")
		lines[4] = sm.IncidentUpdates[0].Format(sm.Shortlink)
	} else {
		for idx, update := range sm.IncidentUpdates {
			url := ""
			if idx == 0 {
				url = sm.Shortlink
			}
			lines[len(sm.IncidentUpdates)-idx] = update.Format(url)
		}
	}
	return strings.Join(lines, "\n")
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
