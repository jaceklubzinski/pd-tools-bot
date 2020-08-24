package team

import (
	"fmt"

	"github.com/jaceklubzinski/pd-tools-bot/pkg/client"
)

type NewTeam struct {
	Team client.TeamClient
}

// PrintTeams list of all available teams in PagerDuty
func (c *NewTeam) PrintTeams() (strs string) {
	client := c.Team.ListTeams()
	for _, p := range client.Teams {
		strstmp := fmt.Sprintf("ID: %s Name: %s\n", p.APIObject.ID, p.Name)
		strs = strs + strstmp
	}
	return strs
}
