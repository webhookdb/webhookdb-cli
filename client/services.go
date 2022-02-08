package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/types"
)

type ServicesListInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

type ServicesListOutput struct {
	Data []ServiceEntity `json:"items"`
}

func ServicesList(c context.Context, auth Auth, input ServicesListInput) (out ServicesListOutput, err error) {
	err = makeRequest(c, GET, auth, input, &out, "/v1/organizations/%v/services", input.OrgIdentifier)
	return
}
