package go_micro_core

import (
	"github.com/facebookgo/inject"
	"go.uber.org/config"
)

type Provider interface {
	Provide() []*inject.Object
}

type ProvideFactory interface {
	NewProvider(provider config.Provider) Provider
}

type NamedInject interface {
	InjectName() string
}

type ProvideFunc func() []*inject.Object

func (f ProvideFunc) Provide() []*inject.Object {
	return f()
}

func NewProvider(vals ...interface{}) ProvideFunc {
	return func() []*inject.Object {
		var objects []*inject.Object
		for _, val := range vals {
			if val == nil {
				continue
			}

			if v, ok := val.(ProvideFactory); ok {
				if p := v.NewProvider(conf); p == nil {
					continue
				} else {
					objects = append(objects, p.Provide()...)
				}
			}

			if v, ok := val.(Provider); ok {
				objects = append(objects, v.Provide()...)
			} else if v, ok := val.(*inject.Object); ok {
				objects = append(objects, v)
			} else if v, ok := val.(NamedInject); ok {
				objects = append(objects, &inject.Object{Name: v.InjectName(), Value: v})
			} else {
				objects = append(objects, &inject.Object{Value: val})
			}
		}

		return objects
	}
}
