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

type DisplayHeaders []KeyAndName

func (dh DisplayHeaders) Names() []string {
	names := make([]string, len(dh))
	for i, h := range dh {
		names[i] = h.Name
	}
	return names
}

type CollectionResponse map[string]interface{}

func (r CollectionResponse) Message() string {
	return respMessage(r)
}

func (r CollectionResponse) DisplayHeaders() DisplayHeaders {
	return respDisplayHeaders(r)
}

func (r CollectionResponse) Items() []map[string]interface{} {
	raw := r["items"].([]interface{})
	maps := make([]map[string]interface{}, len(raw))
	for i, o := range raw {
		maps[i] = o.(map[string]interface{})
	}
	return maps
}

type SingleResponse map[string]interface{}

func (r SingleResponse) Message() string {
	return respMessage(r)
}

func (r SingleResponse) DisplayHeaders() DisplayHeaders {
	return respDisplayHeaders(r)
}

func (r SingleResponse) Fields() map[string]interface{} {
	result := make(map[string]interface{}, len(r))
	for k, v := range r {
		if k != "display_headers" && v != "message" {
			result[k] = v
		}
	}
	return result
}

func respMessage(cr map[string]interface{}) string {
	msg := cr["message"]
	strmsg, ok := msg.(string)
	if !ok {
		fmt.Println("Expected string message in:", cr)
		return ""
	}
	return strmsg
}

func respDisplayHeaders(cr map[string]interface{}) DisplayHeaders {
	raw := cr["display_headers"].([]interface{})
	result := make(DisplayHeaders, len(raw))
	for i, pair := range raw {
		sl := pair.([]interface{})
		result[i] = KeyAndName{Key: sl[0].(string), Name: sl[1].(string)}
	}
	return result
}
