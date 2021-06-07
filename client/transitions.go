package client

import (
	"context"
)

type TransitionStepInput struct {
	AuthCookie string
	PostUrl string
	Value string `json:"value"`
}

func TransitionStep(c context.Context, input TransitionStepInput) (step Step, err error) {
	resty := RestyFromContext(c)
	resp, err := resty.R().
		SetBody(&input).
		SetError(&ErrorResponse{}).
		SetResult(&step).
		SetHeader("Cookie", input.AuthCookie).
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
