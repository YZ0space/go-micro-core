package config

import "go-micro-core/option"

type PostgresqlConnection struct {
}

func (p *PostgresqlConnection) Stop() {
	// close anything
}

func OpenPG(db *option.Postgresql) *PostgresqlConnection {

	return nil
}
