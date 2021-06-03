package client

import (
	"context"
	"strings"
)

type AuthLoginInput struct {
	Username string `json:"email"`
}

type AuthLoginOutput struct {
	Message string `json:"message"`
}

func AuthLogin(c context.Context, input AuthLoginInput) (out AuthLoginOutput, err error) {
	resty := RestyFromContext(c)
	resp, err := resty.R().
		SetBody(&input).
		SetError(&ErrorResponse{}).
		SetResult(&out).
		Post("/v1/auth")
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}

type AuthOTPInput struct {
	Username string `json:"email"`
	Token    string `json:"token"`
}

type AuthOTPResponseOutput struct {
	DefaultOrgKey string `json:"default_org_key"`
	Message string `json:"message"`
}

type AuthOTPOutput struct {
	AuthCookie string
	DefaultOrgKey string
	Message    string
}

func AuthOTP(c context.Context, input AuthOTPInput) (out AuthOTPOutput, err error) {
	resty := RestyFromContext(c)
	respOut := AuthOTPResponseOutput{}
	resp, err := resty.R().
		SetBody(&input).
		SetError(&ErrorResponse{}).
		SetResult(&respOut).
		Post("/v1/auth/login_otp")
	if err != nil {
		return AuthOTPOutput{}, err
	}
	setCookieHeader := resp.Header().Get("Set-Cookie")
	out.AuthCookie = strings.Split(setCookieHeader, ";")[0]
	out.DefaultOrgKey = respOut.DefaultOrgKey
	out.Message = respOut.Message
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}

type AuthLogoutOutput struct {
	Message string `json:"message"`
}

func AuthLogout(c context.Context) (out AuthLogoutOutput, err error) {
	resty := RestyFromContext(c)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetResult(&out).
		Post("/v1/auth/logout")
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}
