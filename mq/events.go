package mq

import "fmt"

// ExampleEvent is used in tests
type ExampleEvent struct {
	Message string `json:"Message"`
}

func (e ExampleEvent) String() string {
	return fmt.Sprintf("ExampleTestEvent, Message=%s", e.Message)
}

// NewSiteEvent is used between repoctl & hugo-helper, as well as sitemgr & k8s-client, to create resources for a new site
type NewSiteEvent struct {
	Email      string `json:"email"`
	SiteName   string `json:"site_name"`
	ReceivedOn int64  `json:"received_on"`
}

func (e NewSiteEvent) String() string {
	return fmt.Sprintf("NewSiteEvent, Email=%s, SiteName=%s, ReceivedOn=%s", e.Email, e.SiteName, e.ReceivedOn)
}

// DelSiteEvent is used between sitemgr & k8s-client, to delete resources of a site
type DelSiteEvent struct {
	Email      string `json:"email"`
	SiteName   string `json:"site_name"`
	ReceivedOn int64  `json:"received_on"`
}

func (e DelSiteEvent) String() string {
	return fmt.Sprintf("DelSiteEvent, Email=%s, SiteName=%s, ReceivedOn=%s", e.Email, e.SiteName, e.ReceivedOn)
}

// ReGenerateEvent is used between repoctl & hugo-helper, to re-generate a site
type ReGenerateEvent struct {
	Email      string `json:"email"`
	SiteName   string `json:"site_name"`
	Theme      string `json:"theme"`
	ReceivedOn int64  `json:"received_on"`
}
