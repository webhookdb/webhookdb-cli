package client

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/types"
	"strconv"
	"strings"
)

type SubscriptionInfoInput struct {
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"identifier"`
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

func (info SubscriptionInfoOutput) PrintTo() {
	fmt.Println("Organization name: " + info.OrgName)
	fmt.Println("Billing email: " + info.BillingEmail)
	fmt.Println("Integrations used: " + strconv.Itoa(info.IntegrationsUsed))
	fmt.Println("Plan name: " + info.PlanName)
	fmt.Println("Integrations left: " + info.IntegrationsLeftDisplay)
	if strings.TrimSpace(info.SubStatus) != "" {
		fmt.Println("Subscription status: " + info.SubStatus)
	}
}

func SubscriptionInfo(c context.Context, input SubscriptionInfoInput) (SubscriptionInfoOutput, error) {
	resty := RestyFromContext(c)
	var out = SubscriptionInfoOutput{}
	url := fmt.Sprintf("/v1/organizations/%v/subscriptions", input.OrgIdentifier)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetResult(&out).
		SetHeader("Cookie", string(input.AuthCookie)).
		Get(url)
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}

type SubscriptionEditInput struct {
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"identifier"`
}

type SubscriptionEditOutput struct {
	SessionUrl string `json:"url"`
}

func SubscriptionEdit(c context.Context, input SubscriptionEditInput) (out SubscriptionEditOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/organizations/%v/subscriptions/open_portal", input.OrgIdentifier)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetResult(&out).
		SetHeader("Cookie", string(input.AuthCookie)).
		Post(url)
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}
