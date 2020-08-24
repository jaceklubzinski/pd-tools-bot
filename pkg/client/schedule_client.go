package client

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/base"
)

type ScheduleClient interface {
	ListSchedules() *pagerduty.ListSchedulesResponse
	ListSchedulesID() *pagerduty.ListSchedulesResponse
}

var scheduleOpts pagerduty.ListSchedulesOptions

func (c *ApiClient) ListSchedules() *pagerduty.ListSchedulesResponse {
	eps, err := c.client.ListSchedules(scheduleOpts)
	base.CheckErr(err)
	return eps
}

func (c *ApiClient) ListSchedulesID() *pagerduty.ListSchedulesResponse {
	eps, err := c.client.ListSchedules(scheduleOpts)
	base.CheckErr(err)
	return eps
}
