package client

import (
	"context"
	"github.com/webhookdb/webhookdb-cli/types"
)

type SyncTargetCreateInput struct {
	OrgIdentifier         types.OrgIdentifier `json:"-"`
	IntegrationIdentifier string              `json:"service_integration_identifier"`
	ConnectionUrl         string              `json:"connection_url"`
	PageSize              int                 `json:"page_size,omitempty"`
	Period                int                 `json:"period_seconds,omitempty"`
	Schema                string              `json:"schema"`
	Table                 string              `json:"table"`
	SyncTypeSlug          string              `json:"-"`
}

func SyncTargetCreate(c context.Context, auth Auth, input SyncTargetCreateInput) (out types.SingleResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/sync_targets/%v/create", input.OrgIdentifier, input.SyncTypeSlug)
	return
}

type SyncTargetDeleteInput struct {
	OpaqueId      string              `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Confirm       string              `json:"confirm"`
	SyncTypeSlug  string              `json:"-"`
}

func SyncTargetDelete(c context.Context, auth Auth, input SyncTargetDeleteInput) (types.MessageResponse, error) {
	out := types.MessageResponse{}
	err := makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/sync_targets/%v/%v/delete", input.OrgIdentifier, input.SyncTypeSlug, input.OpaqueId)
	return out, err
}

type SyncTargetListInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	SyncTypeSlug  string              `json:"-"`
}

func SyncTargetList(c context.Context, auth Auth, input SyncTargetListInput) (out types.CollectionResponse, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v/sync_targets/%v", input.OrgIdentifier, input.SyncTypeSlug)
	return
}

type SyncTargetUpdateInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	OpaqueId      string              `json:"-"`
	PageSize      int                 `json:"page_size,omitempty"`
	Period        int                 `json:"period_seconds,omitempty"`
	Schema        string              `json:"schema"`
	Table         string              `json:"table"`
	SyncTypeSlug  string              `json:"-"`
}

func SyncTargetUpdate(c context.Context, auth Auth, input SyncTargetUpdateInput) (out types.SingleResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/sync_targets/%v/%v/update", input.OrgIdentifier, input.SyncTypeSlug, input.OpaqueId)
	return
}

type SyncTargetUpdateCredsInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	OpaqueId      string              `json:"-"`
	Username      string              `json:"user"`
	Password      string              `json:"password"`
	SyncTypeSlug  string              `json:"-"`
}

func SyncTargetUpdateCreds(c context.Context, auth Auth, input SyncTargetUpdateCredsInput) (out types.SingleResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/sync_targets/%v/%v/update_credentials", input.OrgIdentifier, input.SyncTypeSlug, input.OpaqueId)
	return
}

type SyncTargetSyncInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	OpaqueId      string              `json:"-"`
	SyncTypeSlug  string              `json:"-"`
}

func SyncTargetSync(c context.Context, auth Auth, input SyncTargetSyncInput) (out types.SingleResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/sync_targets/%v/%v/sync", input.OrgIdentifier, input.SyncTypeSlug, input.OpaqueId)
	return
}
