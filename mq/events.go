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
	UserID     uint   `json:"user_id"`
	SiteID     uint   `json:"site_id"`
	SiteName   string `json:"site_name"`
	Theme      string `json:"theme"`
	ReceivedOn int64  `json:"received_on"`
}

func (e SiteUpdatedEvent) String() string {
	return fmt.Sprintf("SiteUpdatedEvent, UserID=%s, SiteId=%s, Theme=%s, ReceivedOn=%s", e.UserID, e.SiteID, e.Theme, e.ReceivedOn)
}
