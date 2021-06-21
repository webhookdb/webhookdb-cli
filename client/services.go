package client

import (
	"context"
)

type ServicesListInput struct {
	AuthCookie string
}

type ServicesListOutput struct {
	Data []ServiceEntity `json:"items"`
}

func ServicesList(c context.Context, input ServicesListInput) (out ServicesListOutput, err error) {
	resty := RestyFromContext(c)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetResult(&out).
		SetHeader("Cookie", input.AuthCookie).
		Get("/v1/services/")
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}
