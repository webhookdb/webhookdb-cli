package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/types"
)

type WebhookCreateInput struct {
	WebhookSecret string              `json:"webhook_secret"`
	Url           string              `json:"url"`
	OrgIdentifier types.OrgIdentifier `json:"org_identifier"`
	SintOpaqueId  string              `json:"service_integration_opaque_id"`
}

func WebhookCreate(c context.Context, auth Auth, input WebhookCreateInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/webhook_subscriptions/create")
	return
}

type WebhookOpaqueIdInput struct {
	// this is the opaque id of the *webhook subscription*
	OpaqueId string `json:"-"`
}

func WebhookTest(c context.Context, auth Auth, input WebhookOpaqueIdInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/webhook_subscriptions/%v/test", input.OpaqueId)
	return
}

func WebhookDelete(c context.Context, auth Auth, input WebhookOpaqueIdInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/webhook_subscriptions/%v/delete", input.OpaqueId)
	return
}
