package client

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/types"
)

type BackfillInput struct {
	AuthCookie types.AuthCookie `json:"-"`
	OpaqueId   string           `json:"-"`
}

func Backfill(c context.Context, input BackfillInput) (step Step, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/service_integrations/%v/backfill", input.OpaqueId)
	resp, err := resty.R().
		SetBody(&input).
		SetError(&ErrorResponse{}).
		SetResult(&step).
		SetHeader("Cookie", string(input.AuthCookie)).
		Post(url)
	if err != nil {
		return step, err
	}
	if err := CoerceError(resp); err != nil {
		return step, err
	}
	return step, nil
}
