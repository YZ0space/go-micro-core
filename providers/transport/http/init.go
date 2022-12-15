package http

import (
	"github.com/aka-yz/go-micro-core"
)

func init() {
	go_micro_core.RegisterProvider(
		&serverFactory{},
	)
}
