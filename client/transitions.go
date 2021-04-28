package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/statemachine"
)

type TransitionStepInput struct {
	Input string `json:"input"`
}

type TransitionStepOutput struct {
	Step statemachine.Step `json:"step"`
}

func TransitionStep(c context.Context, input TransitionStepInput) (out TransitionStepOutput, err error) {
	resty := RestyFromContext(c)
	resp, err := resty.R().
		SetBody(&input).
		SetError(&ErrorResponse{}).
		SetResult(&out).
		Post("/v1/state_machine/transition_step")
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}

func NewStateMachine() statemachine.StateMachine {
	return statemachine.New(func(c context.Context, input string) (statemachine.Step, error) {
		output, err := TransitionStep(c, TransitionStepInput{Input: input})
		if err != nil {
			return statemachine.Step{}, err
		}
		return output.Step, nil
	})
}
