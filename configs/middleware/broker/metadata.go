package broker

import "context"
import "strings"

type metadataKey struct{}

type Metadata map[string][]string

func MetadataFromContext(ctx context.Context) Metadata {
	return ctx.Value(metadataKey{}).(Metadata)
}

func NewContextFromMetadata(meta Metadata) context.Context {
	return context.WithValue(context.Background(), metadataKey{}, meta)
}

func (m Metadata) Set(key, val string) {
	key = strings.ToLower(key)
	m[key] = []string{val}
}

func (m Metadata) Get(key string) string {
	key = strings.ToLower(key)
	var v string
	if vv, ok := m[key]; ok && len(vv) > 0 {
		v = vv[0]
	}
	return v
}

func Join(mds ...Metadata) Metadata {
	out := Metadata{}
	for _, md := range mds {
		for k, v := range md {
			out[k] = append(out[k], v...)
		}
	}
	return out
}

func NewMetadata(m map[string][]string, prefix string) Metadata {
	md := Metadata{}
	prefix = strings.ToLower(prefix)
	for k, v := range m {
		k = strings.ToLower(k)
		if strings.HasPrefix(k, prefix+"-") {
			md[k] = v
		}
	}
	return md
}
