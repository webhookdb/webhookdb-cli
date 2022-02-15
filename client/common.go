package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lithictech/webhookdb-cli/types"
	"net/http"
)

const RestyKey = "client.resty"

func RestyInContext(c context.Context, r *resty.Client) context.Context {
	return context.WithValue(c, RestyKey, r)
}

func RestyFromContext(c context.Context) *resty.Client {
	return c.Value(RestyKey).(*resty.Client)
}

type Auth struct {
	Cookie types.AuthCookie `json:"-"`
}

type ErrorResponse struct {
	Err struct {
		Message          string `json:"message"`
		Code             string `json:"code"`
		Status           int    `json:"status"`
		StateMachineStep *Step  `json:"state_machine_step"`
	} `json:"error"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("ErrorResponse[%d, %s]: %s", e.Err.Status, e.Err.Code, e.Err.Message)
}

func CoerceError(r *resty.Response) error {
	if r.StatusCode() >= 400 {
		return *r.Error().(*ErrorResponse)
	}
	return nil
}

const (
	GET  = http.MethodGet
	POST = http.MethodPost
)

func makeRequest(c context.Context, method string, auth Auth, body, outPtr interface{}, urlTmpl string, urlArgs ...interface{}) error {
	resp, errStep, err := makeRequestWithResponse(c, method, auth, body, outPtr, urlTmpl, urlArgs...)
	if err != nil && err != errStepNoInputAndNotComplete {
		// If we have an error and it's not just an error because
		// the state machine finished and we got a real response, return it.
		return err
	}
	if resp != nil {
		// The call worked, so we can assume the output ptr got bound properly.
		return nil
	}
	// The state machine completed, but this was a 'normal' endpoint, so we have to parse the final response body
	// back to the requested body (the final response body is not a step, while the errors were).
	if err := json.Unmarshal(errStep.RawResponse.Body(), outPtr); err != nil {
		return err
	}
	return nil
}

func makeStepRequestWithResponse(c context.Context, auth Auth, body interface{}, urlTmpl string, urlArgs ...interface{}) (Step, error) {
	var step Step
	resp, errStep, err := makeRequestWithResponse(c, POST, auth, body, &step, urlTmpl, urlArgs...)
	if err != nil {
		return step, err
	}
	// This is safe to reprocess, it will early-out if it came from an error step
	// since we know step endpoints that have a terminated step are finished for good.
	if errStep != nil {
		return *errStep, nil
	}
	// This is the 'true' state machine step (no intercepting error state machine)
	// so we need to set the RawResponse.
	step.RawResponse = resp
	return step, err
}

// Return the response from the given request, the Step from a 426 error output, and any error.
// Either the response, OR the error step, will be returned on 'success'.
// They will not both be returned; to get the raw response of the error step,
// use the Step.RawResponse field.
func makeRequestWithResponse(c context.Context, method string, auth Auth, body, outPtr interface{}, urlTmpl string, urlArgs ...interface{}) (*resty.Response, *Step, error) {
	r := RestyFromContext(c)
	url := fmt.Sprintf(urlTmpl, urlArgs...)
	req := r.R().SetError(&ErrorResponse{})
	if body != nil {
		req = req.SetBody(body)
	}
	if outPtr != nil {
		req = req.SetResult(outPtr)
	}
	if auth.Cookie != "" {
		req = req.SetHeader("Cookie", string(auth.Cookie))
	}
	resp, err := req.Execute(method, url)
	if err != nil {
		return resp, nil, err
	}
	if err := CoerceError(resp); err != nil {
		if eresp, ok := err.(ErrorResponse); ok && eresp.Err.StateMachineStep != nil {
			errStep, err := NewStateMachine().RunWithOutput(c, auth, *eresp.Err.StateMachineStep)
			errStep.processedViaError = true
			return nil, &errStep, err
		}
		return resp, nil, err
	}
	return resp, nil, nil
}
