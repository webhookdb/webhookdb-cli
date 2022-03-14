package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/types"
)

func AuthGetMe(c context.Context, auth Auth) (out types.SingleResponse, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/me")
	return
}

type AuthLoginInput struct {
	Username string `json:"email"`
	Token    string `json:"token"`
}

type AuthLoginOutput struct {
	AuthToken  types.AuthToken
	OutputStep Step
}

func AuthLogin(c context.Context, input AuthLoginInput) (Step, error) {
	return makeStepRequestWithResponse(c, Auth{}, input, "/v1/auth")
}

type AuthLogoutOutput struct {
	Message string `json:"message"`
}

func AuthLogout(c context.Context, auth Auth) (out AuthLogoutOutput, err error) {
	err = makeRequest(c, POST, auth, nil, &out, "/v1/auth/logout")
	return
}
