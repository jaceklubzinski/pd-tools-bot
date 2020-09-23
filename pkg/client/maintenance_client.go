package client

import (
	"github.com/PagerDuty/go-pagerduty"
)

type MaintenanceClient interface {
	ListMaintenance(opts pagerduty.ListMaintenanceWindowsOptions) (*pagerduty.ListMaintenanceWindowsResponse, error)
	CreateMaintenance(o pagerduty.MaintenanceWindow) (*pagerduty.MaintenanceWindow, error)
}

func (c *APIClient) ListMaintenance(opts pagerduty.ListMaintenanceWindowsOptions) (*pagerduty.ListMaintenanceWindowsResponse, error) {
	eps, err := c.client.ListMaintenanceWindows(opts)
	return eps, err
}

func (c *APIClient) CreateMaintenance(o pagerduty.MaintenanceWindow) (*pagerduty.MaintenanceWindow, error) {
	from := "pdbot"
	eps, err := c.client.CreateMaintenanceWindow(from, o)
	return eps, err
}
