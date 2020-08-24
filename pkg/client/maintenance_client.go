package client

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/base"
)

type MaintenanceClient interface {
	ListMaintenance(opts pagerduty.ListMaintenanceWindowsOptions) *pagerduty.ListMaintenanceWindowsResponse
	CreateMaintenance(o pagerduty.MaintenanceWindow) *pagerduty.MaintenanceWindow
}

func (c *ApiClient) ListMaintenance(opts pagerduty.ListMaintenanceWindowsOptions) *pagerduty.ListMaintenanceWindowsResponse {
	eps, err := c.client.ListMaintenanceWindows(opts)
	base.CheckErr(err)
	return eps
}

func (c *ApiClient) CreateMaintenance(o pagerduty.MaintenanceWindow) *pagerduty.MaintenanceWindow {
	from := "pdbot"
	eps, err := c.client.CreateMaintenanceWindow(from, o)
	base.CheckErr(err)
	return eps
}
