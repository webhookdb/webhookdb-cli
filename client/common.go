package client

import (
	"context"
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
		Message string `json:"message"`
		Code    string `json:"code"`
		Status  int    `json:"status"`
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
	GET    = http.MethodGet
	POST   = http.MethodPost
	PUT    = http.MethodPut
	DELETE = http.MethodDelete
)

func makeRequest(c context.Context, method string, auth Auth, body, outPtr interface{}, urlTmpl string, urlArgs ...interface{}) error {
	_, err := makeRequestWithResponse(c, method, auth, body, outPtr, urlTmpl, urlArgs...)
	return err
}

func makeStepRequestWithResponse(c context.Context, method string, auth Auth, body interface{}, urlTmpl string, urlArgs ...interface{}) (Step, error) {
	var step Step
	resp, err := makeRequestWithResponse(c, method, auth, body, &step, urlTmpl, urlArgs...)
	if err != nil {
		return step, err
	}
	step.RawResponse = resp
	return step, err
}

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
	if auth.Cookie != "" {
		req = req.SetHeader("Cookie", string(auth.Cookie))
	}
	resp, err := req.Execute(method, url)
	if err != nil {
		return resp, err
	}
	if err := CoerceError(resp); err != nil {
		return resp, err
	}
	return resp, nil
}
