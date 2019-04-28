package db

import (
	"database/sql"
)

type Miners struct {
	Sub_account_name sql.NullString `db:"sub_account_name"`
	Sub_account_id   int64          `db:"sub_account_id"`
}

func (c *Client) QuerySubAccount(last_id string, limit int) ([]*Miners, error) {
	queryStr :=
		`SELECT
	f.sub_account_name,	
	f.sub_account_id
FROM
	f_sub_account f
WHERE f.deltag = '0'
and  sub_account_id > ?
order by sub_account_id
limit ?`

	var miners []*Miners
	err := c.client.Select(&miners, queryStr, last_id, limit)
	return miners, err
}
