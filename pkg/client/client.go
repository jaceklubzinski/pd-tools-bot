package client

import (
	"github.com/PagerDuty/go-pagerduty"
)

type ApiClient struct {
	client *pagerduty.Client
}

func NewApiClient(client *pagerduty.Client) *ApiClient {
	return &ApiClient{client: client}
}
