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
	Token types.AuthToken `json:"-"`
}

const AuthTokenHeader = "Whdb-Auth-Token"

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
	_, err := makeRequestWithResponse(c, method, auth, body, outPtr, urlTmpl, urlArgs...)
	return err
}

func makeStepRequestWithResponse(c context.Context, auth Auth, body interface{}, urlTmpl string, urlArgs ...interface{}) (Step, error) {
	var step Step
	resp, err := makeRequestWithResponse(c, POST, auth, body, &step, urlTmpl, urlArgs...)
	step.RawResponse = resp
	return step, err
}

// Return the response from the given request, the Step from a 422 prompt, and any error.
// Either the response, OR the error step, will be returned on 'success'.
// They will not both be returned; to get the raw response of the error step,
// use the Step.RawResponse field.
// Note that, since the 422 prompt handling is recursive, if this request 422s,
// but the recursive state machine eventually succeeds, we will return that successful result.
func makeRequestWithResponse(c context.Context, method string, auth Auth, body, outPtr interface{}, urlTmpl string, urlArgs ...interface{}) (*resty.Response, error) {
	r := RestyFromContext(c)
	url := fmt.Sprintf(urlTmpl, urlArgs...)
	req := r.R().SetError(&ErrorResponse{})
	if body != nil {
		req = req.SetBody(body)
	}
	if outPtr != nil {
		req = req.SetResult(outPtr)
	}
	if auth.Token != "" {
		req = req.SetHeader(AuthTokenHeader, string(auth.Token))
	}
	resp, err := req.Execute(method, url)
	if err != nil {
		return resp, err
	}
	if err := CoerceError(resp); err != nil {
		if eresp, ok := err.(ErrorResponse); ok && eresp.Err.StateMachineStep != nil {
			errStep, err := NewStateMachine().RunWithOutput(c, auth, *eresp.Err.StateMachineStep)
			if err != nil {
				return errStep.RawResponse, err
			}
			if err := json.Unmarshal(errStep.RawResponse.Body(), outPtr); err != nil {
				return errStep.RawResponse, err
			}
			return errStep.RawResponse, nil
		}
		return resp, err
	}
	return resp, nil
}
