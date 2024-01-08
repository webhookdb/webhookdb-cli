package client

import (
	"context"
	"github.com/webhookdb/webhookdb-cli/types"
)

type ServicesListInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func ServicesList(c context.Context, auth Auth, input ServicesListInput) (out types.CollectionResponse, err error) {
	err = makeRequest(c, GET, auth, input, &out, "/v1/organizations/%v/services", input.OrgIdentifier)
	return
}
