package client

import (
	"github.com/PagerDuty/go-pagerduty"
)

type IncidentClient interface {
	ListIncidents(opts pagerduty.ListIncidentsOptions) (*pagerduty.ListIncidentsResponse, error)
}

func (c *APIClient) ListIncidents(opts pagerduty.ListIncidentsOptions) (*pagerduty.ListIncidentsResponse, error) {
	eps, err := c.client.ListIncidents(opts)
	return eps, err
}
