package bot

import (
	"strings"

	"github.com/gowee/github-status-bot/pkg/api"
)

func formatMultipleComponents(components []api.Component) string {
	seq := make([]string, len(components))
	for idx, component := range components {
		if component.ID == "0l2p9nhqnxpd" {
			// skip one entry
		}
		seq[idx] = component.Format()
	}
	return strings.Join(seq, "\n")
}
