package client

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/types"
	"io"
	"strings"
)

type AuthCurrentCustomerOutput struct {
	Message             string                         `json:"message"`
	Email               string                         `json:"email"`
	Name                string                         `json:"name"`
	DefaultOrganization types.Organization             `json:"default_organization"`
	Memberships         []OrganizationMembershipEntity `json:"memberships"`
}

func (o AuthCurrentCustomerOutput) PrintTo(w io.Writer) {
	fmt.Fprintln(w, "Name: "+o.Name)
	fmt.Fprintln(w, "Email: "+o.Email)
	fmt.Fprintln(w, "Default Org: "+o.DefaultOrganization.DisplayString())
	if len(o.Memberships) == 0 {
		fmt.Fprintln(w, "Memberships: <none>")
	} else {
		fmt.Fprintln(w, "Memberships:")
		for _, m := range o.Memberships {
			fmt.Fprintf(w, "  %s: %s\n", m.Organization.DisplayString(), m.Status)
		}
	}
}

func AuthGetMe(c context.Context, auth Auth) (out AuthCurrentCustomerOutput, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/me")
	return
}

type AuthLoginInput struct {
	Username string `json:"email"`
}

type AuthLoginOutput struct {
	Message string `json:"message"`
}

func AuthLogin(c context.Context, input AuthLoginInput) (out AuthLoginOutput, err error) {
	err = makeRequest(c, POST, Auth{}, input, &out, "/v1/auth")
	return
}

type AuthOTPInput struct {
	Username string `json:"email"`
	Token    string `json:"token"`
}

type AuthOTPOutput struct {
	Message         string
	AuthCookie      types.AuthCookie
	CurrentCustomer AuthCurrentCustomerOutput
}

func AuthOTP(c context.Context, input AuthOTPInput) (out AuthOTPOutput, err error) {
	resty := RestyFromContext(c)
	respOut := AuthCurrentCustomerOutput{}
	resp, err := resty.R().
		SetBody(&input).
		SetError(&ErrorResponse{}).
		SetResult(&respOut).
		Post("/v1/auth/login_otp")
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	setCookieHeader := resp.Header().Get("Set-Cookie")
	out.AuthCookie = types.AuthCookie(strings.Split(setCookieHeader, ";")[0])
	out.CurrentCustomer = respOut
	out.Message = respOut.Message
	return out, nil
}

type AuthLogoutOutput struct {
	Message string `json:"message"`
}

func AuthLogout(c context.Context, auth Auth) (out AuthLogoutOutput, err error) {
	err = makeRequest(c, POST, auth, nil, &out, "/v1/auth/logout")
	return
}
