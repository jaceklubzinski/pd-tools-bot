package schedule

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/client"
)

type NewSchedule struct {
	Schedule client.ScheduleClient
}

// PrintSchedules list of all available schedules in PagerDuty
func (c *NewSchedule) PrintSchedules() (strs string) {
	client := c.Schedule.ListSchedules()
	for _, p := range client.Schedules {
		strstmp := fmt.Sprintf("ID: %s Name: %s\n", p.APIObject.ID, p.Name)
		strs = strs + strstmp
	}
	return strs
}

// PrintSchedules list of all available schedules in PagerDuty
func (c *NewSchedule) GetAll() (schedules []pagerduty.Schedule) {
	client := c.Schedule.ListSchedules()
	return client.Schedules
}

// GetScheduleID transform schedule name to ID
func (c *NewSchedule) GetScheduleID(scheduleName string) (scheduleID string) {
	client := c.Schedule.ListSchedulesID()
	for _, p := range client.Schedules {
		if p.Name == scheduleName {
			scheduleID = p.APIObject.ID
		}
	}
	return
}
