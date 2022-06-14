package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/types"
)

type WebhookCreateInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	WebhookSecret string              `json:"webhook_secret"`
	Url           string              `json:"url"`
	SintOpaqueId  string              `json:"service_integration_opaque_id"`
}

func WebhookCreate(c context.Context, auth Auth, input WebhookCreateInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/webhook_subscriptions/create", input.OrgIdentifier)
	return
}

type WebhookListInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func WebhookList(c context.Context, auth Auth, input WebhookListInput) (out types.CollectionResponse, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v/webhook_subscriptions", input.OrgIdentifier)
	return
}

type WebhookOpaqueIdInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	// this is the opaque id of the *webhook subscription*
	OpaqueId string `json:"-"`
}

func WebhookTest(c context.Context, auth Auth, input WebhookOpaqueIdInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/webhook_subscriptions/%v/test", input.OrgIdentifier, input.OpaqueId)
	return
}

func WebhookDelete(c context.Context, auth Auth, input WebhookOpaqueIdInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/webhook_subscriptions/%v/delete", input.OrgIdentifier, input.OpaqueId)
	return
}
