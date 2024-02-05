package client

import (
	"context"
	"github.com/webhookdb/webhookdb-cli/formatting"
	"github.com/webhookdb/webhookdb-cli/types"
)

type SavedQueryListInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func SavedQueryList(c context.Context, auth Auth, input SavedQueryListInput) (out types.CollectionResponse, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v/saved_queries", input.OrgIdentifier)
	return
}

type SavedQueryCreateInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Description   string              `json:"description"`
	Sql           string              `json:"sql"`
	Public        bool                `json:"public"`
}

func SavedQueryCreate(c context.Context, auth Auth, input SavedQueryCreateInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%s/saved_queries/create", input.OrgIdentifier)
	return
}

type SavedQueryIdentifierInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Identifier    string              `json:"-"`
}

func SavedQueryDelete(c context.Context, auth Auth, input SavedQueryIdentifierInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/saved_queries/%v/delete", input.OrgIdentifier, input.Identifier)
	return
}

type SavedQueryUpdateInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Identifier    string              `json:"-"`
	Field         string              `json:"field"`
	Value         string              `json:"value"`
}

func SavedQueryUpdate(c context.Context, auth Auth, input SavedQueryUpdateInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/saved_queries/%v/update", input.OrgIdentifier, input.Identifier)
	return
}

type SavedQueryInfoInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Identifier    string              `json:"-"`
	Field         string              `json:"field"`
}

type SavedQueryInfoOutput struct {
	Blocks formatting.Blocks `json:"blocks"`
}

func SavedQueryInfo(c context.Context, auth Auth, input SavedQueryInfoInput) (out SavedQueryInfoOutput, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/saved_queries/%v/info", input.OrgIdentifier, input.Identifier)
	return
}

func SavedQueryRun(c context.Context, auth Auth, input SavedQueryIdentifierInput) (out types.RunQueryOutput, err error) {
	err = makeRequest(c, GET, auth, input, &out, "/v1/organizations/%v/saved_queries/%v/run", input.OrgIdentifier, input.Identifier)
	return
}
