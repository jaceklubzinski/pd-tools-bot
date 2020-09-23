package team

import (
	"fmt"

	"github.com/jaceklubzinski/pd-tools-bot/pkg/client"
)

//Team PagerDuty client for team
type Team struct {
	Team client.TeamClient
}

// PrintTeams list of all available teams in PagerDuty
func (c *Team) PrintTeams() (strs string) {
	client := c.Team.ListTeams()
	for _, p := range client.Teams {
		strstmp := fmt.Sprintf("ID: %s Name: %s\n", p.APIObject.ID, p.Name)
		strs = strs + strstmp
	}
	return strs
}
