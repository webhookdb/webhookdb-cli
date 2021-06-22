package client

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/types"
)

type DbTablesInput struct {
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

type DbTablesOutput struct {
	TableNames []string `json:"tables"`
}

func DbTables(c context.Context, input DbTablesInput) (out DbTablesOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/db/%v", input.OrgIdentifier)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetResult(&out).
		SetHeader("Cookie", string(input.AuthCookie)).
		Get(url)
	fmt.Println(resp)
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
	fmt.Println(resp)
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}
