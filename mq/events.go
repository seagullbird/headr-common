package mq

import "fmt"

// ExampleEvent is used in tests
type ExampleEvent struct {
	Message string `json:"Message"`
}

func (e ExampleEvent) String() string {
	return fmt.Sprintf("ExampleTestEvent, Message=%s", e.Message)
}

// SiteUpdatedEvent is used between repoctl & hugo-helper, as well as sitemgr & k8s-client, to generate site
type SiteUpdatedEvent struct {
	SiteId     uint   `json:"site_id"`
	Theme      string `json:"theme"`
	ReceivedOn int64  `json:"received_on"`
}

func (e SiteUpdatedEvent) String() string {
	return fmt.Sprintf("SiteUpdatedEvent, SiteId=%s, Theme=%s, ReceivedOn=%s", e.SiteId, e.Theme, e.ReceivedOn)
}
