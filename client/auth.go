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
	Token    string `json:"token"`
}

type AuthLoginOutput struct {
	AuthCookie types.AuthCookie
	OutputStep Step
}

func AuthLogin(c context.Context, input AuthLoginInput) (out AuthLoginOutput, err error) {
	resty := RestyFromContext(c)
	step := Step{}
	resp, err := resty.R().
		SetBody(&input).
		SetError(&ErrorResponse{}).
		SetResult(&step).
		Post("/v1/auth")
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	setCookieHeader := resp.Header().Get("Set-Cookie")
	out.AuthCookie = types.AuthCookie(strings.Split(setCookieHeader, ";")[0])
	out.OutputStep = step
	return out, nil
}

type AuthLogoutOutput struct {
	Message string `json:"message"`
}

func AuthLogout(c context.Context, auth Auth) (out AuthLogoutOutput, err error) {
	err = makeRequest(c, POST, auth, nil, &out, "/v1/auth/logout")
	return
}
