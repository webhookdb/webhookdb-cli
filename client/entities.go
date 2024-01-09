package client

import "github.com/webhookdb/webhookdb-cli/types"

type OrganizationMembershipEntity struct {
	CustomerEmail string             `json:"email"`
	Organization  types.Organization `json:"organization"`
	Status        string             `json:"status"`
}

type ServiceEntity struct {
	Name string `json:"name"`
}

type ServiceIntegrationEntity struct {
	OpaqueId    string `json:"opaque_id"`
	ServiceName string `json:"service_name"`
	TableName   string `json:"table_name"`
}

type WebhookSubscriptionEntity struct {
	OpaqueId           string                   `json:"opaque_id"`
	DeliverToUrl       string                   `json:"deliver_to_url"`
	Organization       types.Organization       `json:"organization"`
	ServiceIntegration ServiceIntegrationEntity `json:"service_integration"`
}
