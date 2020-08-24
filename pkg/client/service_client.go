package client

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/base"
)

type ServiceClient interface {
	ListServices(opts pagerduty.ListServiceOptions) *pagerduty.ListServiceResponse
}

func (c *ApiClient) ListServices(opts pagerduty.ListServiceOptions) *pagerduty.ListServiceResponse {
	eps, err := c.client.ListServices(opts)
	base.CheckErr(err)
	return eps
}
