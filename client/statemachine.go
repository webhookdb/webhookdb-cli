package client

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/ask"
	"github.com/lithictech/webhookdb-cli/prefs"
)

type Step struct {
	Message        string `json:"message"`
	NeedsInput     bool   `json:"needs_input"`
	Prompt         string `json:"prompt"`
	PromptIsSecret bool   `json:"prompt_is_secret"`
	PostToUrl      string `json:"post_to_url"`
	Complete       bool   `json:"complete"`
	Output         string `json:"output"`
}

type Prompt func(string) (string, error)
type Println func(...interface{})
//type TransitionStep func(c context.Context, input TransitionStepInput) (Step, error)
//
func NewStateMachine() StateMachine {
	return StateMachine{
		Ask:            ask.Ask,
		HiddenAsk:      ask.HiddenAsk,
		Println:        func(a ...interface{}) { fmt.Println(a...) },
		//TransitionStep: transition,
	}
}

type StateMachine struct {
	Ask       Prompt
	HiddenAsk Prompt
	Println   Println
	//TransitionStep TransitionStep
}

func (sm StateMachine) Run(c context.Context, p prefs.Prefs, startingStep Step) error {
	step := startingStep
	for {
		if step.Complete {
			sm.Println(step.Output)
			return nil
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
			return err
		}
		transitionInput := TransitionStepInput{
			AuthCookie: p.AuthCookie,
			PostUrl: step.PostToUrl,
			Value:      value,
		}
		newStep, err := TransitionStep(c, transitionInput)
		if err != nil {
			return err
		}
		step = newStep
	}
}