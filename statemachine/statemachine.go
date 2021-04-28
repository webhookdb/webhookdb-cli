package statemachine

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/ask"
)

type Step struct {
	NeedsInput     bool   `json:"needs_input"`
	Prompt         string `json:"prompt"`
	PromptIsSecret bool   `json:"prompt_is_secret"`
	PostToUrl      string `json:"post_to_url"`
	Complete       bool   `json:"complete"`
	Output         string `json:"output"`
}

type Prompt func(string) (string, error)
type Println func(...interface{})
type TransitionStep func(c context.Context, input string) (Step, error)

func New(transition TransitionStep) StateMachine {
	return StateMachine{
		Ask:            ask.Ask,
		HiddenAsk:      ask.HiddenAsk,
		Println:        func(a ...interface{}) { fmt.Println(a...) },
		TransitionStep: transition,
	}
}

type StateMachine struct {
	Ask            Prompt
	HiddenAsk      Prompt
	Println        Println
	TransitionStep TransitionStep
}

func (sm StateMachine) Run(c context.Context, startingStep Step) error {
	step := startingStep
	for {
		if step.Complete {
			sm.Println(step.Output)
			return nil
		}
		if !step.NeedsInput {
			panic("Step must be complete, or need input. Backend is busted.")
		}
		asker := sm.Ask
		prompt := step.Prompt
		if step.PromptIsSecret {
			asker = sm.HiddenAsk
			prompt = ask.HiddenPrompt(prompt)
		}
		input, err := asker(step.Prompt)
		if err != nil {
			return err
		}
		newStep, err := sm.TransitionStep(c, input)
		if err != nil {
			return err
		}
		step = newStep
	}
}
