package client

import (
	"github.com/PagerDuty/go-pagerduty"
)

type APIClient struct {
	client *pagerduty.Client
}

func NewAPIClient(client *pagerduty.Client) *APIClient {
	return &APIClient{client: client}
}
