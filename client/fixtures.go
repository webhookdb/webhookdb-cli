package client

import "context"

type GetFixturesInput struct {
	ServiceName string `json:"service_name"`
}

type GetFixturesOutput struct {
	SchemaSql string `json:"schema_sql"`
}

func GetFixtures(c context.Context, auth Auth, input GetFixturesInput) (out GetFixturesOutput, err error) {
	err = makeRequest(c, GET, auth, input, &out, "/v1/services/%v/fixtures", input.ServiceName)
	return
}
