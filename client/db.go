package client

import (
	"context"
	"encoding/json"
	"github.com/lithictech/webhookdb-cli/types"
)

type DbOrgIdentifierInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

type DbConnectionOutput struct {
	ConnectionUrl string `json:"connection_url"`
}

func DbConnection(c context.Context, auth Auth, input DbOrgIdentifierInput) (out DbConnectionOutput, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/db/%v/connection", input.OrgIdentifier)
	return
}

type DbTablesOutput struct {
	Message    string   `json:"message"`
	TableNames []string `json:"tables"`
}

func DbTables(c context.Context, auth Auth, input DbOrgIdentifierInput) (out DbTablesOutput, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/db/%v/tables", input.OrgIdentifier)
	return
}

type DbSqlInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Query         string              `json:"query"`
}

type DbSqlOutput struct {
	// Use RawMessage to avoid deserializing the JSON right away.
	// This allows us to render maps and certain other types verbatim.
	Rows           [][]json.RawMessage `json:"rows"`
	Headers        []string            `json:"headers"`
	MaxRowsReached bool                `json:"max_rows_reached"`
}

func DbSql(c context.Context, auth Auth, input DbSqlInput) (out DbSqlOutput, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/db/%v/sql", input.OrgIdentifier)
	return
}

func DbRollCredentials(c context.Context, auth Auth, input DbOrgIdentifierInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, nil, &out, "/v1/db/%v/roll_credentials", input.OrgIdentifier)
	return
}

type DbFdwInput struct {
	OrgIdentifier    types.OrgIdentifier `json:"-"`
	MessageFdw       bool                `json:"message_fdw"`
	MessageViews     bool                `json:"message_views"`
	MessageAll       bool                `json:"message_all"`
	RemoteServerName string              `json:"remote_server_name"`
	FetchSize        string              `json:"fetch_size"`
	LocalSchema      string              `json:"local_schema"`
	ViewSchema       string              `json:"view_schema"`
}

type DbFdwOutput map[string]interface{}

func DbFdw(c context.Context, auth Auth, input DbFdwInput) (out DbFdwOutput, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/db/%v/fdw", input.OrgIdentifier)
	return
}

type DbRenameTableInput struct {
	IntegrationIdentifier string              `json:"-"`
	OrgIdentifier         types.OrgIdentifier `json:"-"`
	NewName               string              `json:"new_name"`
}

func DbRenameTable(c context.Context, auth Auth, input DbRenameTableInput) (out types.SingleResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/service_integrations/%v/rename_table", input.OrgIdentifier, input.IntegrationIdentifier)
	return
}

type DbMigrationsStartInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	AdminUrl      string              `json:"admin_url"`
	ReadonlyUrl   *string             `json:"readonly_url,omitempty"`
}

func DbMigrationsStart(c context.Context, auth Auth, input DbMigrationsStartInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/db/%v/migrate_database", input.OrgIdentifier)
	return
}

func DbMigrationsList(c context.Context, auth Auth, input DbOrgIdentifierInput) (out types.CollectionResponse, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/db/%v/migrations", input.OrgIdentifier)
	return
}
