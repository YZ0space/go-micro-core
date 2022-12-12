package config

import (
	"go-micro-core/option"
)

type DBConnection struct {
	//*dbr.Connection
}

func (p *DBConnection) Stop() {
	// close anything
}

func OpenDB(db *option.DB) *DBConnection {
	return nil
}
