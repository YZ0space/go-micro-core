package db

import (
	"github.com/gocraft/dbr/v2/dialect"
	"time"
)

type postgres struct{}

func (d postgres) QuoteIdent(s string) string {
	return dialect.PostgreSQL.QuoteIdent(s)
}

func (d postgres) EncodeString(s string) string {
	return dialect.PostgreSQL.EncodeString(s)
}

func (d postgres) EncodeBool(b bool) string {
	return dialect.PostgreSQL.EncodeBool(b)
}

func (d postgres) EncodeTime(t time.Time) string {
	return `'` + t.Format(timeFormat) + `'`
}

func (d postgres) EncodeBytes(b []byte) string {
	return dialect.PostgreSQL.EncodeBytes(b)
}

func (d postgres) Placeholder(n int) string {
	return dialect.PostgreSQL.Placeholder(n)
}
