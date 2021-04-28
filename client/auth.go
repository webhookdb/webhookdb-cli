package client

import (
	"context"
)

type AuthRegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRegisterOutput struct {
}

func AuthRegister(c context.Context, input AuthRegisterInput) (out AuthRegisterOutput, err error) {
	resty := RestyFromContext(c)
	resp, err := resty.R().
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
