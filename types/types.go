package types

import (
	"fmt"
	"strconv"
)

type OrgIdentifier string

func OrgIdentifierFromId(id int) OrgIdentifier {
	return OrgIdentifier(strconv.Itoa(id))
}

func OrgIdentifierFromSlug(slug string) OrgIdentifier {
	return OrgIdentifier(slug)
}

type AuthCookie string

type Organization struct {
	// The primary key of the organization.
	// Safe to store on the client.
	Id int `json:"id" mapstructure:"id"`
	// The name of the org. Can change.
	Name string `json:"name" mapstructure:"name"`
	// URL-safe version of the organization name.
	Key string `json:"key" mapstructure:"key"`
}

func (o Organization) DisplayString() string {
	if o.Id == 0 {
		return "<none>"
	}
	return fmt.Sprintf("%s (%s)", o.Name, o.Key)
}
