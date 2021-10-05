package client

import (
	"context"
)

type TransitionStepInput struct {
	PostUrl string `json:"-"`
	Value   string `json:"value"`
}

func TransitionStep(c context.Context, auth Auth, input TransitionStepInput) (step Step, err error) {
	err = makeRequest(c, POST, auth, input, &step, input.PostUrl)
	return
}
