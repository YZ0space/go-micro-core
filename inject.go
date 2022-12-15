package go_micro_core

import (
	"github.com/aka-yz/go-micro-core/extension"
	"github.com/facebookgo/inject"
)

type Inject struct {
	inject.Graph
	Vals []interface{}
}

var injects Inject

func RegisterProvider(objs ...interface{}) {
	injects.Vals = append(injects.Vals, objs...)
}

func ScanGinHandler(handlerName string) extension.GinHandler {
	for _, v := range injects.Objects() {
		if o, ok := v.Value.(extension.GinHandler); ok {
			if v.Name == handlerName {
				return o
			}
		}
	}
	return nil
}
