package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type StatusPage struct {
	PageID string
}

// QueryOverall querys QuerySummary along with QueryIncidents and QueryScheduledMaintanances to
// fill up all incidents and scheduled maintenances instead of only active ones as is the case for
// QuerySummary.
func (sp *StatusPage) QueryOverall() (qres *SummaryQueryResult, err error) {
	summary, err := sp.QuerySummary()
	if err != nil {
		return nil, err
	}
	allIncidents, err := sp.QueryIncidents()
	if err != nil {
		return nil, err
	}
	summary.Incidents = allIncidents.Incidents
	allScheduledMaintenances, err := sp.QueryScheduledMaintanances()
	if err != nil {
		return nil, err
	}
	summary.ScheduledMaintenances = allScheduledMaintenances.ScheduledMaintenances

	return &summary, nil
}

func (sp *StatusPage) QuerySummary() (qres SummaryQueryResult, err error) {
	url := sp.urlFor("summary")
	qres = SummaryQueryResult{}
	err = getJSON(url, &qres)
	return
}

func (sp *StatusPage) QueryIncidents() (qres IncidentsQueryResult, err error) {
	url := sp.urlFor("incidents")
	qres = IncidentsQueryResult{}
	err = getJSON(url, &qres)
	return
}

func (sp *StatusPage) QueryScheduledMaintanances() (qres ScheduledMaintenancesQueryResult, err error) {
	url := sp.urlFor("scheduled-maintenances")
	qres = ScheduledMaintenancesQueryResult{}
	err = getJSON(url, &qres)
	return
}

func (sp *StatusPage) urlFor(endpoint string) string {
	return fmt.Sprintf("https://%s.statuspage.io/api/v2/%s.json", sp.PageID, endpoint)
}

func getJSON(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}

// https://kctbh9vrtdwd.statuspage.io/api/v2/status.json
