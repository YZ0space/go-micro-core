package db

import (
	"context"
	"fmt"
	"github.com/aka-yz/go-micro-core/providers/option"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Connection struct {
	*pgx.Conn
}

func (d *Connection) Stop() {
	if err := d.Close(context.Background()); err != nil {
		panic(err)
	}
}

//func (d *Connection) NewSession() *dbr.Session {
//	return d.Connection.NewSession(nil)
//}

func OpenDB(option *option.DB) *Connection {
	if option.Port == 0 {
		option.Port = 5432
	}

	if option.Host == "" {
		option.Host = "localhost"
	}

	// postgres://username:password@localhost:5432/database_name
	if option.Driver == "postgres" {
		option.DataSource = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			option.UserName,
			option.Password,
			option.Host,
			option.Port,
			option.DBName,
		)
	}

	var conn *pgx.Conn
	var err error
	//var logEventReceiver = NewEventReceiver(dbName(option.DataSource), 200, 200)
	conn, err = pgx.Connect(context.Background(), option.DataSource)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to connect to database: %v\n", err)
		panic(err)
	}
	return &Connection{Conn: conn}
}
