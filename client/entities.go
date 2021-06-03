package client

type OrganizationEntity struct {
	Key string `json:"key"`
}

type OrganizationMembershipEntity struct {
	CustomerEmail string `json:"email"`
	Status string `json:"status"`
}

type ServiceEntity struct {
	Name string `json:"name"`
}
