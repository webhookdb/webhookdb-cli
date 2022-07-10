package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/types"
)

type NamedQueryCreateInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Name          string              `json:"name"`
	Sql           string              `json:"sql"`
}

type NamedQueryCreateOutput struct {
	NamedQueryEntity
	types.MessageResponse
}

func NamedQueryCreate(c context.Context, auth Auth, input NamedQueryCreateInput) (out NamedQueryCreateOutput, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/custom_queries/create", input.OrgIdentifier)
	return
}

type NamedQueryListInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func NamedQueryList(c context.Context, auth Auth, input NamedQueryListInput) (out types.CollectionResponse, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v/custom_queries", input.OrgIdentifier)
	return
}

type NamedQueryLookupInput struct {
	OrgIdentifier   types.OrgIdentifier `json:"-"`
	QueryIdentifier string              `json:"query_identifier"`
}

func NamedQueryInfo(c context.Context, auth Auth, input NamedQueryLookupInput) (out types.SingleResponse, err error) {
	err = makeRequest(c, GET, auth, input, &out, "/v1/organizations/%v/custom_queries/lookup?query_identifier=%s", input.OrgIdentifier, input.QueryIdentifier)
	return
}

func NamedQueryRun(c context.Context, auth Auth, input NamedQueryLookupInput) (out DbSqlOutput, err error) {
	err = makeRequest(c, GET, auth, input, &out, "/v1/organizations/%v/custom_queries/run?query_identifier=%s", input.OrgIdentifier, input.QueryIdentifier)
	return
}

type NamedQueryUpdateInput struct {
	OrgIdentifier   types.OrgIdentifier `json:"-"`
	QueryIdentifier string              `json:"query_identifier"`
	Field           string              `json:"field"`
	Value           string              `json:"value"`
}

func NamedQueryUpdate(c context.Context, auth Auth, input NamedQueryUpdateInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/custom_queries/update", input.OrgIdentifier)
	return
}

func NamedQueryDelete(c context.Context, auth Auth, input NamedQueryLookupInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/custom_queries/delete", input.OrgIdentifier)
	return
}
