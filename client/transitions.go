package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/types"
)

type TransitionStepInput struct {
	AuthCookie types.AuthCookie `json:"-"`
	PostUrl    string           `json:"-"`
	Value      string           `json:"value"`
}

func TransitionStep(c context.Context, input TransitionStepInput) (step Step, err error) {
	resty := RestyFromContext(c)
	resp, err := resty.R().
		SetBody(&input).
		SetError(&ErrorResponse{}).
		SetResult(&step).
		SetHeader("Cookie", string(input.AuthCookie)).
		Post(input.PostUrl)
	if err != nil {
		return step, err
	}
	if err := CoerceError(resp); err != nil {
		return step, err
	}
	return step, nil
}

//func NewStateMachine() statemachine.StateMachine {
//	return statemachine.New(func(c context.Context, input string) (statemachine.Step, error) {
//		output, err := TransitionStep(c, TransitionStepInput{Input: input})
//		if err != nil {
//			return statemachine.Step{}, err
//		}
//		return output.Step, nil
//	})
//}
