package client

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/base"
)

type TeamClient interface {
	ListTeams() *pagerduty.ListTeamResponse
}

var teamOpts pagerduty.ListTeamOptions

func (c *APIClient) ListTeams() *pagerduty.ListTeamResponse {
	eps, err := c.client.ListTeams(teamOpts)
	base.CheckErr(err)
	return eps
}
