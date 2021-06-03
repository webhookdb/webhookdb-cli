package client

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const RestyKey = "client.resty"

func RestyInContext(c context.Context, r *resty.Client) context.Context {
	return context.WithValue(c, RestyKey, r)
}

func RestyFromContext(c context.Context) *resty.Client {
	return c.Value(RestyKey).(*resty.Client)
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
