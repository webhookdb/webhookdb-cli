package client

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/types"
)

type SubscriptionInfoInput struct {
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"identifier"`
}

type SubscriptionInfoOutput struct {
	OrgName          string `json:"org_name"`
	BillingEmail     string `json:"billing_email"`
	IntegrationsUsed string `json:"integrations_used"`
	PlanName         string `json:"plan_name"`
	IntegrationsLeft string `json:"integrations_left"`
	SubStatus        string `json:"sub_status"`
}

func SubscriptionInfo(c context.Context, input SubscriptionInfoInput) (out SubscriptionInfoOutput, err error) {
	resty := RestyFromContext(c)
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
