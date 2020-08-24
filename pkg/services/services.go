package services

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/client"
)

type Services struct {
	Service client.ServiceClient
}

func (s *Services) GetAll() {
	var opts pagerduty.ListServiceOptions
	getAll := s.Service.ListServices(opts)
	for _, p := range getAll.Services {
		fmt.Printf("ID: %s Name: %s Team: %s Ack Timeout: %d\n", p.APIObject.ID, p.Name, p.Teams, p.AcknowledgementTimeout)
	}
}

func (s *Services) GetTeam(teamID []string) (strs string) {
	opts := pagerduty.ListServiceOptions{
		TeamIDs: teamID,
	}
	getTeam := s.Service.ListServices(opts)
	for _, p := range getTeam.Services {
		strstmp := fmt.Sprintf("ID: %s Name: %s Ack Timeout: %d\n", p.APIObject.ID, p.Name, p.AcknowledgementTimeout)
		strs = strs + strstmp
	}
	return strs
}
