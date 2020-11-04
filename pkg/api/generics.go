package api

import "time"

// An Event is either a Incident or a ScheduledMaintenance
type Event interface {
	GetID() string
	GetUpdatedAt() time.Time
	Format() string
	ShouldNotify() bool
	IsFinished() bool
}

func (i *Incident) GetID() string {
	return i.ID
}

func (i *Incident) GetUpdatedAt() time.Time {
	return i.UpdatedAt
}

func (sm *ScheduledMaintenance) GetID() string {
	return sm.ID
}

func (sm *ScheduledMaintenance) GetUpdatedAt() time.Time {
	return sm.UpdatedAt
}

type Events []Event

func (es Events) Len() int {
	return len(es)
}

func (es Events) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es Events) Less(i, j int) bool {
	return es[i].GetUpdatedAt().Before(es[j].GetUpdatedAt())
}
