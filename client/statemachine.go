package client

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lithictech/webhookdb-cli/ask"
)

type Step struct {
	Message        string           `json:"message"`
	NeedsInput     bool             `json:"needs_input"`
	Prompt         string           `json:"prompt"`
	PromptIsSecret bool             `json:"prompt_is_secret"`
	PostToUrl      string           `json:"post_to_url"`
	Complete       bool             `json:"complete"`
	Output         string           `json:"output"`
	Extras         map[string]Extra `json:"extras"`
	RawResponse    *resty.Response  `json:"-"`
}

type Extra map[string]interface{}

type Prompt func(string) (string, error)
type Println func(...interface{})

func NewStateMachine() StateMachine {
	return StateMachine{
		Ask:       ask.Ask,
		HiddenAsk: ask.HiddenAsk,
		Println:   func(a ...interface{}) { fmt.Println(a...) },
	}
}

type StateMachine struct {
	Ask       Prompt
	HiddenAsk Prompt
	Println   Println
}

func (sm StateMachine) Run(c context.Context, auth Auth, startingStep Step) error {
	_, err := sm.RunWithOutput(c, auth, startingStep)
	return err
}

func (sm StateMachine) RunWithOutput(c context.Context, auth Auth, startingStep Step) (Step, error) {
	step := startingStep
	for {
		if step.Complete {
			sm.Println(step.Output)
			return step, nil
		}
		if !step.NeedsInput {
			panic("Step must be complete, or need input. Backend is busted.")
		}
		sm.Println(step.Output)
		asker := sm.Ask
		prompt := step.Prompt
		if step.PromptIsSecret {
			asker = sm.HiddenAsk
			prompt = ask.HiddenPrompt(prompt)
		}
		value, err := asker(step.Prompt)
		if err != nil {
			return step, err
		}
		transitionInput := TransitionStepInput{
			PostUrl: step.PostToUrl,
			Value:   value,
		}
		newStep, err := TransitionStep(c, auth, transitionInput)
		if err != nil {
			return newStep, err
		}
		step = newStep
	}
}
