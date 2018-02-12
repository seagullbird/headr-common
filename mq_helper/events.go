package mq_helper

import "fmt"

type NewSiteEvent struct {
	Email 		string	`json:"email"`
	SiteName 	string	`json:"site_name"`
	ReceivedOn	int64	`json:"received_on"`
}

func (e NewSiteEvent) String() string {
	return fmt.Sprintf("Email: %s, Site Name: %s, Received On: %s", e.Email, e.SiteName, e.ReceivedOn)
}
