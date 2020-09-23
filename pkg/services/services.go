package services

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/base"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/client"
)

//Services Pagerduty client for services
type Services struct {
	Service client.ServiceClient
}

//GetAll list all PagerDuty services
func (s *Services) GetAll() {
	var opts pagerduty.ListServiceOptions
	getAll, err := s.Service.ListServices(opts)
	base.CheckErr(err)
	for _, p := range getAll.Services {
		fmt.Printf("ID: %s Name: %s Team: %s Ack Timeout: %d\n", p.APIObject.ID, p.Name, p.Teams, p.AcknowledgementTimeout)
	}
}

//GetTeam list all PagerDuty services for specific teams
func (s *Services) GetTeam(teamID []string) (strs string, err error) {
	opts := pagerduty.ListServiceOptions{
		TeamIDs: teamID,
	}
	getTeam, err := s.Service.ListServices(opts)
	if err != nil {
		return "", err
	}
	for _, p := range getTeam.Services {
		strstmp := fmt.Sprintf("ID: %s Name: %s Ack Timeout: %d\n", p.APIObject.ID, p.Name, p.AcknowledgementTimeout)
		strs = strs + strstmp
	}
	return strs, nil
}
