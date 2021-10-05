package client

import (
	"context"
)

type ServicesListInput struct{}

type ServicesListOutput struct {
	Data []ServiceEntity `json:"items"`
}

func ServicesList(c context.Context, auth Auth, input ServicesListInput) (out ServicesListOutput, err error) {
	err = makeRequest(c, GET, auth, input, &out, "/v1/services")
	return
}
