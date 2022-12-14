package db

import (
	"fmt"
	"github.com/aka-yz/go-micro-core/providers/option"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"time"
)

type Connection struct {
	*dbr.Connection
}

func (d *Connection) Stop() {
	if err := d.Close(); err != nil {
		panic(err)
	}
}

func (d *Connection) NewSession() *dbr.Session {
	return d.Connection.NewSession(nil)
}

func OpenDB(option *option.DB) *Connection {
	if option.Port == 0 {
		option.Port = 3306
	}

	if option.Host == "" {
		option.Host = "localhost"
	}

	if option.DataSource == "" {
		option.DataSource = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			option.UserName,
			option.Password,
			fmt.Sprintf("%s:%d", option.Host, option.Port),
			option.DBName,
		)
		timeout := option.Timeout
		if timeout == 0 {
			timeout = 3
		}
		option.DataSource += fmt.Sprintf("&timeout=%ds", timeout)

		readTimeout := option.ReadTimeout
		if readTimeout == 0 {
			readTimeout = 8
		}
		option.DataSource += fmt.Sprintf("&readTimeout=%ds", readTimeout)

		writeTimeout := option.WriteTimeout
		if writeTimeout == 0 {
			writeTimeout = 8
		}
		option.DataSource += fmt.Sprintf("&writeTimeout=%ds", writeTimeout)
	}
	if option.Driver == "" {
		option.Driver = "mysql"
	}

	var conn *dbr.Connection
	var err error
	var logEventReceiver = NewEventReceiver(dbName(option.DataSource), 200, 200)
	conn, err = dbr.Open(option.Driver, option.DataSource, logEventReceiver)
	if err != nil {
		panic(err)
	}
	if option.MaxIdleConns == 0 {
		option.MaxIdleConns = 200
	}
	if option.MaxOpenConns == 0 {
		option.MaxOpenConns = 200
	}
	if option.ConnMaxLifetime == 0 {
		option.ConnMaxLifetime = 600
	}
	conn.Dialect = &mysql{}
	conn.SetMaxIdleConns(option.MaxIdleConns)
	conn.SetMaxOpenConns(option.MaxOpenConns)
	conn.SetConnMaxLifetime(time.Duration(option.ConnMaxLifetime) * time.Second)
	return &Connection{Connection: conn}
}
