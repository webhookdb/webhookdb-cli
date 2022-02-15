package client

import (
	"context"
)

type TransitionStepInput struct {
	PostUrl            string
	PostParams         map[string]interface{}
	PostParamsValueKey string
	Value              string
}

func TransitionStep(c context.Context, auth Auth, input TransitionStepInput) (Step, error) {
	params := input.PostParams
	if params == nil {
		params = map[string]interface{}{}
	}
	params[input.PostParamsValueKey] = input.Value
	return makeStepRequestWithResponse(c, auth, params, input.PostUrl)
}
