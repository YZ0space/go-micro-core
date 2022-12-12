package go_micro_core

import (
	"os"
	"strings"
)

type ENV string

var (
	Env ENV
)

func (e ENV) Live() bool {
	return e == "live"
}

func (e ENV) Dev() bool {
	return e == "dev" || e == ""
}

func initEnv() {
	Env = ENV(strings.ToLower(os.Getenv("env")))
}
