package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lithictech/webhookdb-cli/ask"
)

var errStepNoInputAndNotComplete = errors.New("step is not complete, and says it needs no input. The endpoint is configuring the step wrong, or a normal response was parsed as a Step")

type Step struct {
	Message            string                 `json:"message"`
	NeedsInput         bool                   `json:"needs_input"`
	Prompt             string                 `json:"prompt"`
	PromptIsSecret     bool                   `json:"prompt_is_secret"`
	PostToUrl          string                 `json:"post_to_url"`
	PostParams         map[string]interface{} `json:"post_params"`
	PostParamsValueKey string                 `json:"post_params_value_key"`
	Complete           bool                   `json:"complete"`
	Output             string                 `json:"output"`
	Extras             map[string]Extra       `json:"extras"`
	RawResponse        *resty.Response        `json:"-"`
	// If true, this step was processed via the automatic error handling state machine
	// and clients should generally ignore it.
	processedViaError bool
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
	if startingStep.processedViaError {
		// We don't want to run this sort of step multiple times.
		return startingStep, nil
	}
	step := startingStep
	for {
		if step.Complete {
			sm.Println(step.Output)
			return step, nil
		}
		if !step.NeedsInput {
			return step, errStepNoInputAndNotComplete
		}
		if step.Output != "" {
			// If the step is the first one, so only prompts, this will be blank.
			sm.Println(step.Output)
		}
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
			PostUrl:            step.PostToUrl,
			PostParams:         step.PostParams,
			PostParamsValueKey: step.PostParamsValueKey,
			Value:              value,
		}
		newStep, err := TransitionStep(c, auth, transitionInput)
		if err != nil {
			return newStep, err
		}
		step = newStep
		// Always print a newline after processing input, so the next step output
		// has a blank line after the input.
		fmt.Println("")
	}
}

// StateMachineResponseRunner is a helper to wrap client calls that return (Step, error)
// so we can use a single line to make the API call and run the state machine.
func StateMachineResponseRunner(ctx context.Context, auth Auth) func(Step, error) (Step, error) {
	return func(step Step, err error) (Step, error) {
		if err != nil {
			return step, err
		}
		if err := NewStateMachine().Run(ctx, auth, step); err != nil {
			return step, err
		}
		return step, nil
	}
}
