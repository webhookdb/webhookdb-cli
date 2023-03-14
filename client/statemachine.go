package client

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/lithictech/webhookdb-cli/ask"
)

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
}

type Extra map[string]interface{}

func NewStateMachine() StateMachine {
	return StateMachine{
		ask: ask.New(),
	}
}

type StateMachine struct {
	ask ask.Ask
}

// Run processes the input step until it succeeds or the processing errors (from the API, or ctrl+c, etc).
func (sm StateMachine) Run(c context.Context, auth Auth, startingStep Step) (Step, error) {
	if startingStep.Complete {
		// If we pass in a complete step, assume the caller took car of printing the output.
		// Otherwise we can end up wish a finished step from a 422 statemachine, and then re-print the result.
		return startingStep, nil
	}
	step := startingStep
	for {
		if step.Complete {
			sm.ask.Feedback(step.Output)
			return step, nil
		}
		if !step.NeedsInput {
			// Usually this is because a 422 prompt machine returned success
			return step, nil
		}
		if step.Output != "" {
			// If the step is the first one, so only prompts, this will be blank.
			sm.ask.Feedback(step.Output)
		}
		asker := sm.ask.Ask
		prompt := step.Prompt
		if step.PromptIsSecret {
			asker = sm.ask.HiddenAsk
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
		if eresp, ok := err.(ErrorResponse); ok && eresp.Err.Code == "validation_error" {
			// If the field that fails validation is not the one we submitted, it probably means that
			// something from the commandline failed. There's no sense re-prompting if the current field is valid,
			// since we can't fix the cause of the 400 through this transition.
			if _, newValueIsInvalid := eresp.Err.FieldErrors[step.PostParamsValueKey]; newValueIsInvalid {
				// Print the message and offer the same prompt again for new input.
				sm.ask.Feedback(eresp.Err.Message)
				sm.ask.Feedback("")
				continue
			} else {
				return newStep, err
			}
		} else if err != nil {
			return newStep, err
		}
		step = newStep
		// Always print a newline after processing input, so the next step output
		// has a blank line after the input.
		sm.ask.Feedback("")
	}
}

// StateMachineResponseRunner is a helper to wrap client calls that return (Step, error)
// so we can use a single line to make the API call and run the state machine.
func StateMachineResponseRunner(ctx context.Context, auth Auth) func(Step, error) (Step, error) {
	return func(step Step, err error) (Step, error) {
		if err != nil {
			return step, err
		}
		newStep, newErr := NewStateMachine().Run(ctx, auth, step)
		if newErr != nil {
			return step, err
		}
		return newStep, nil
	}
}
