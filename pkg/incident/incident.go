package incident

import (
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/base"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/client"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/extensions"
)

//Incidents PagerDuty client
type Incidents struct {
	Incident client.IncidentClient
}

//GetAll list all PagerDuty incidents
func (i *Incidents) GetAll() {
	var opts pagerduty.ListIncidentsOptions
	getAll, err := i.Incident.ListIncidents(opts)
	base.CheckErr(err)
	for _, p := range getAll.Incidents {
		fmt.Printf("ID: %s Name: %s  Service: %s\n", p.APIObject.ID, p.Title, p.Service.Summary)
	}
}

//GetTeam list PagerDuty incidents for specifc teams
func (i *Incidents) GetTeam(teamID []string) (strs string, err error) {
	serviceIncidents := make(map[string]string)
	serviceIncidentsNumber := make(map[string]int)
	opts := pagerduty.ListIncidentsOptions{
		APIListObject: pagerduty.APIListObject{
			Limit: 100,
		},
		TeamIDs:  teamID,
		SortBy:   "created_at",
		Statuses: []string{"triggered", "acknowledged"},
	}
	getTeam, err := i.Incident.ListIncidents(opts)
	if err != nil {
		return "", err
	}
	for _, p := range getTeam.Incidents {
		strstmp := fmt.Sprintf("Name: %s Created At: %s Status: %s <%s|PD Link>\n", p.Title, p.CreatedAt, p.Status, p.HTMLURL)
		serviceIncidents[p.Service.Summary] = serviceIncidents[p.Service.Summary] + strstmp
		serviceIncidentsNumber[p.Service.Summary]++
	}
	for k, v := range serviceIncidents {
		strstmp := fmt.Sprintf("\n`%s - %d`\n%s", k, serviceIncidentsNumber[k], v)
		strs = strs + strstmp
	}
	return strs, nil
}

//GetTeamDuty list PagerDuty incidents summary after duty
func (i *Incidents) GetTeamDuty(teamID []string, startHour string) (strs string, err error) {
	serviceIncidents := make(map[string]string)
	serviceIncidentsNumber := make(map[string]int)
	startToday := time.Now().Format("2006-01-02 15:04:05")
	sinceDate, err := extensions.BackDurationToDate(startToday, startHour)
	if err != nil {
		return "", err
	}
	opts := pagerduty.ListIncidentsOptions{
		APIListObject: pagerduty.APIListObject{
			Limit: 100,
		},
		TeamIDs:  teamID,
		Since:    sinceDate,
		Until:    startToday,
		SortBy:   "created_at",
		Statuses: []string{"triggered", "acknowledged", "resolved"},
	}
	getTeam, err := i.Incident.ListIncidents(opts)
	if err != nil {
		return "", err
	}
	for _, p := range getTeam.Incidents {
		strstmp := fmt.Sprintf("Name: %s Created At: %s Status: %s <%s|PD Link>\n", p.Title, p.CreatedAt, p.Status, p.HTMLURL)
		serviceIncidents[p.Service.Summary] = serviceIncidents[p.Service.Summary] + strstmp
		serviceIncidentsNumber[p.Service.Summary]++
	}
	for k, v := range serviceIncidents {
		strstmp := fmt.Sprintf("\n`%s - %d`\n%s", k, serviceIncidentsNumber[k], v)
		strs = strs + strstmp
	}
	return strs, nil
}
