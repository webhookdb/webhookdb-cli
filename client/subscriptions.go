package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/types"
)

type SubscriptionInfoInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func SubscriptionInfo(c context.Context, auth Auth, input SubscriptionInfoInput) (out types.SingleResponse, err error) {
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
}

func SubscriptionPlans(c context.Context, auth Auth, input SubscriptionPlansInput) (out types.CollectionResponse, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v/subscriptions/plans", input.OrgIdentifier)
	return
}
