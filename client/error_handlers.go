package client

import (
	"context"
	"github.com/webhookdb/webhookdb-cli/types"
)

type ErrorHandlerListInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func ErrorHandlerList(c context.Context, auth Auth, input ErrorHandlerListInput) (out types.CollectionResponse, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v/error_handlers", input.OrgIdentifier)
	return
}

type ErrorHandlerCreateInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Url           string              `json:"url"`
}

func ErrorHandlerCreate(c context.Context, auth Auth, input ErrorHandlerCreateInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%s/error_handlers/create", input.OrgIdentifier)
	return
}

type ErrorHandlerIdentifierInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Identifier    string              `json:"-"`
}

func ErrorHandlerDelete(c context.Context, auth Auth, input ErrorHandlerIdentifierInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/error_handlers/%v/delete", input.OrgIdentifier, input.Identifier)
	return
}
