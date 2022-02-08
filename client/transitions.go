package client

import (
	"context"
)

type TransitionStepInput struct {
	PostUrl string `json:"-"`
	Value   string `json:"value"`
}

func TransitionStep(c context.Context, auth Auth, input TransitionStepInput) (Step, error) {
	return makeStepRequestWithResponse(c, POST, auth, input, input.PostUrl)
}
