package client

import (
	"context"
	"github.com/webhookdb/webhookdb-cli/types"
)

type SavedViewListInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func SavedViewList(c context.Context, auth Auth, input SavedViewListInput) (out types.CollectionResponse, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v/saved_views", input.OrgIdentifier)
	return
}

type SavedViewCreateInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Name          string              `json:"name"`
	Sql           string              `json:"sql"`
}

func SavedViewCreate(c context.Context, auth Auth, input SavedViewCreateInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%s/saved_views/create_or_replace", input.OrgIdentifier)
	return
}

type SavedViewDeleteInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Name          string              `json:"name"`
}

func SavedViewDelete(c context.Context, auth Auth, input SavedViewDeleteInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/saved_views/delete", input.OrgIdentifier)
	return
}
