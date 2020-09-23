package client

import (
	"github.com/PagerDuty/go-pagerduty"
)

type OnCallClient interface {
	ListOnCalls(opts pagerduty.ListOnCallOptions) (*pagerduty.ListOnCallsResponse, error)
}

// ListUsers on calls users
func (c *APIClient) ListOnCalls(opts pagerduty.ListOnCallOptions) (*pagerduty.ListOnCallsResponse, error) {
	eps, err := c.client.ListOnCalls(opts)
	return eps, err
}
