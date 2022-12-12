package go_micro_core

import "github.com/facebookgo/inject"

type Inject struct {
	inject.Graph
	Vals []interface{}
}

var injects Inject

func RegisterProvider(objs ...interface{}) {
	injects.Vals = append(injects.Vals, objs...)
}
