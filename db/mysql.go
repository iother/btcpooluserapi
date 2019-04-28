package db

import (
	"github.com/go-sql-driver/mysql"
	"github.com/iother/btcpooluserapi/config"
	"github.com/iother/btcpooluserapi/log"
	"github.com/jmoiron/sqlx"
)

type Client struct {
	client *sqlx.DB
}

func NewMySQL(config *config.DB) (*Client, error) {
	conf := &mysql.Config{
		User:                 config.User,
		Passwd:               config.Password,
		Addr:                 config.Addr,
		Net:                  "tcp",
		DBName:               config.Database,
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	db, err := sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Error(" mysql Open fail ", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Error(" mysql Ping fail ", err)
		return nil, err
	}
	return &Client{db}, nil
}

func (c *Client) CloseDB() {
	c.client.Close()
}
