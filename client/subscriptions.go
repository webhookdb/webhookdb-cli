package client

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/formatting"
	"github.com/lithictech/webhookdb-cli/types"
	"io"
	"strconv"
	"strings"
)

type SubscriptionInfoInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

type SubscriptionInfoOutput struct {
	OrgName                 string `json:"org_name"`
	BillingEmail            string `json:"billing_email"`
	IntegrationsUsed        int    `json:"integrations_used"`
	PlanName                string `json:"plan_name"`
	IntegrationsLeft        int    `json:"integrations_left"`
	IntegrationsLeftDisplay string `json:"integrations_left_display"`
	SubStatus               string `json:"sub_status"`
}

func (info SubscriptionInfoOutput) PrintTo(w io.Writer) {
	fmt.Fprintln(w, "Organization name: "+info.OrgName)
	fmt.Fprintln(w, "Billing email: "+info.BillingEmail)
	fmt.Fprintln(w, "Integrations used: "+strconv.Itoa(info.IntegrationsUsed))
	fmt.Fprintln(w, "Plan name: "+info.PlanName)
	fmt.Fprintln(w, "Integrations left: "+info.IntegrationsLeftDisplay)
	if strings.TrimSpace(info.SubStatus) != "" {
		fmt.Fprintln(w, "Subscription status: "+info.SubStatus)
	}
}

func SubscriptionInfo(c context.Context, auth Auth, input SubscriptionInfoInput) (out SubscriptionInfoOutput, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v/subscriptions", input.OrgIdentifier)
	return
}

type SubscriptionEditInput struct {
	OrgIdentifier types.OrgIdentifier `json:"identifier"`
	Plan          string              `json:"plan"`
}

type SubscriptionEditOutput struct {
	SessionUrl string `json:"url"`
}

func SubscriptionEdit(c context.Context, auth Auth, input SubscriptionEditInput) (out SubscriptionEditOutput, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/subscriptions/open_portal", input.OrgIdentifier)
	return
}

type SubscriptionPlansInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Format        formatting.Format   `json:"-"`
}

type SubscriptionPlansOutput struct {
	Parsed interface{}
}

func SubscriptionPlans(c context.Context, auth Auth, input SubscriptionPlansInput) (SubscriptionPlansOutput, error) {
	out := SubscriptionPlansOutput{
		Parsed: input.Format.ApiResponsePtr(),
	}
	err := makeRequest(c, GET, auth, nil, out.Parsed, "/v1/organizations/%v/subscriptions/plans?fmt=%s", input.OrgIdentifier, input.Format.ApiRequestValue)
	return out, err
}
