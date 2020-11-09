package client

import (
	"log"

	"github.com/PagerDuty/go-pagerduty"
)

type IncidentClient interface {
	ListIncidents(opts pagerduty.ListIncidentsOptions) (*pagerduty.ListIncidentsResponse, error)
	ListIncidentNotes(id string) (string, error)
}

func (c *APIClient) ListIncidents(opts pagerduty.ListIncidentsOptions) (*pagerduty.ListIncidentsResponse, error) {
	eps, err := c.client.ListIncidents(opts)
	return eps, err
}

func (c *APIClient) ListIncidentNotes(id string) (string, error) {
	var notes string
	eps, err := c.client.ListIncidentNotes(id)
	if err != nil {
		log.Fatalln(err)
	}
	for _, n := range eps {
		notes = notes + n.Content
	}
	return notes, err
}
