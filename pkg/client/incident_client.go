package client

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/base"
)

type IncidentClient interface {
	ListIncidents(opts pagerduty.ListIncidentsOptions) *pagerduty.ListIncidentsResponse
}

func (c *ApiClient) ListIncidents(opts pagerduty.ListIncidentsOptions) *pagerduty.ListIncidentsResponse {
	eps, err := c.client.ListIncidents(opts)
	base.CheckErr(err)
	return eps
}
