package client

import (
	"github.com/PagerDuty/go-pagerduty"
)

type ServiceClient interface {
	ListServices(opts pagerduty.ListServiceOptions) (*pagerduty.ListServiceResponse, error)
}

func (c *APIClient) ListServices(opts pagerduty.ListServiceOptions) (*pagerduty.ListServiceResponse, error) {
	eps, err := c.client.ListServices(opts)
	return eps, err
}
