package middleware

import (
	"github.com/aka-yz/go-micro-core/providers/option"
	"github.com/go-pg/pg/v10"
	"time"
)

type PostgresqlConnection struct {
	*pg.DB
}

func (p *PostgresqlConnection) Stop() {
	// close anything
	if err := p.Close(); err != nil {
		panic(err)
	}
}

func OpenPG(opt *option.Postgresql) *PostgresqlConnection {
	db := pg.Connect(&pg.Options{
		Addr:                  opt.Addr,
		User:                  opt.User,
		Password:              opt.Password,
		Database:              opt.Database,
		ApplicationName:       opt.ApplicationName,
		DialTimeout:           0,
		ReadTimeout:           50000,
		WriteTimeout:          50000,
		MaxRetries:            5,
		RetryStatementTimeout: false,
		PoolSize:              100,
		MinIdleConns:          5,
		MaxConnAge:            2 * time.Hour,
		PoolTimeout:           1 * time.Hour,
		IdleTimeout:           2 * time.Hour,
		IdleCheckFrequency:    10 * time.Minute,
	})
	return &PostgresqlConnection{db}
}
