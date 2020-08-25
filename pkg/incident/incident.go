package incident

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/base"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/client"
)

type Incidents struct {
	Incident client.IncidentClient
}

func (i *Incidents) GetAll() {
	var opts pagerduty.ListIncidentsOptions
	getAll, err := i.Incident.ListIncidents(opts)
	base.CheckErr(err)
	for _, p := range getAll.Incidents {
		fmt.Printf("ID: %s Name: %s  Service: %s\n", p.APIObject.ID, p.Title, p.Service.Summary)
	}
}

func (i *Incidents) GetTeam(teamID []string) (strs string, err error) {
	opts := pagerduty.ListIncidentsOptions{
		TeamIDs:  teamID,
		Statuses: []string{"triggered", "acknowledged"},
	}
	getTeam, err := i.Incident.ListIncidents(opts)
	if err != nil {
		return "", err
	}
	for _, p := range getTeam.Incidents {
		strstmp := fmt.Sprintf("ID: %s Name: %s  Service: %s Status: %s\n", p.APIObject.ID, p.Title, p.Service.Summary, p.Status)
		strs = strs + strstmp
	}
	return strs, nil
}
