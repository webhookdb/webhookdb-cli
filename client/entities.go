package client

type OrganizationMembershipEntity struct {
	CustomerEmail string `json:"email"`
	Status        string `json:"status"`
}

type ServiceEntity struct {
	Name string `json:"name"`
}

type ServiceIntegrationEntity struct {
	OpaqueId    string `json:"opaque_id"`
	ServiceName string `json:"service_name"`
	TableName   string `json:"table_name"`
}
