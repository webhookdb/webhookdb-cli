package client

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/types"
)

type DbConnectionInput struct {
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

type DbConnectionOutput struct {
	ConnectionUrl string `json:"connection_url"`
}

func DbConnection(c context.Context, input DbConnectionInput) (out DbConnectionOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/db/%v/connection", input.OrgIdentifier)
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

type DbTablesInput struct {
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

type DbTablesOutput struct {
	TableNames []string `json:"tables"`
}

func DbTables(c context.Context, input DbTablesInput) (out DbTablesOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/db/%v/tables", input.OrgIdentifier)
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

type DbSqlInput struct {
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Query         string              `json:"query"`
}

type DbSqlOutput struct {
	Rows           []string `json:"rows"`
	Columns        []string `json:"columns"`
	MaxRowsReached bool     `json:"max_rows_reached"`
}

func DbSql(c context.Context, input DbSqlInput) (out DbSqlOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/db/%v/sql", input.OrgIdentifier)
	resp, err := resty.R().
		SetBody(input).
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

type DbRollCredentialsInput struct {
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

type DbRollCredentialsOutput struct {
	ConnectionUrl string `json:"connection_url"`
}

func DbRollCredentials(c context.Context, input DbRollCredentialsInput) (out DbRollCredentialsOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/db/%v/roll_credentials", input.OrgIdentifier)
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
