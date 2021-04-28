package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
)

type AuthRegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRegisterOutput struct {
}

func AuthRegister(c context.Context, input AuthRegisterInput) (out AuthRegisterOutput, err error) {
	actx := appcontext.FromContext(c)
	resp, err := actx.Resty.R().
		SetBody(&input).
		SetError(&ErrorResponse{}).
		SetResult(&out).
		Post("/v1/auth/register")
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}
