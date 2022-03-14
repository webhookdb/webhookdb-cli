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

type AuthToken string

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

type MessageResponse struct {
	Message string `json:"message"`
}

type KeyAndName struct{ Key, Name string }

type CollectionResponse map[string]interface{}

func (cr CollectionResponse) Message() string {
	msg := cr["message"]
	strmsg, ok := msg.(string)
	if !ok {
		fmt.Println("Expected string message in:", cr)
		return ""
	}
	return strmsg
}

func (cr CollectionResponse) Items() []map[string]interface{} {
	raw := cr["items"].([]interface{})
	maps := make([]map[string]interface{}, len(raw))
	for i, o := range raw {
		maps[i] = o.(map[string]interface{})
	}
	return maps
}

func (cr CollectionResponse) DisplayHeaders() []KeyAndName {
	raw := cr["display_headers"].([]interface{})
	result := make([]KeyAndName, len(raw))
	for i, pair := range raw {
		sl := pair.([]interface{})
		result[i] = KeyAndName{Key: sl[0].(string), Name: sl[1].(string)}
	}
	return result
}
