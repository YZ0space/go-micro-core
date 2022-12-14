package metadata

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"strings"
)

type MD map[string][]string

func New(m map[string]string) MD {
	md := MD{}
	for k, val := range m {
		key := strings.ToLower(k)
		md[key] = append(md[key], val)
	}
	return md
}

func Pairs(kv ...string) MD {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: Pairs got the odd number of input pairs for metadata: %d", len(kv)))
	}
	md := MD{}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = strings.ToLower(s)
			continue
		}
		md[key] = append(md[key], s)
	}
	return md
}

// Len returns the number of items in md.
func (md MD) Len() int {
	return len(md)
}

// Copy returns a copy of md.
func (md MD) Copy() MD {
	return Join(md)
}

// Get obtains the values for a given key.
func (md MD) Get(k string) string {
	k = strings.ToLower(k)
	var v string
	if vv, ok := md[k]; ok && len(vv) > 0 {
		v = vv[0]
	}
	return v
}

// Set sets the value of a given key with a slice of values.
func (md MD) Set(k string, val string) {
	if len(val) == 0 {
		return
	}
	k = strings.ToLower(k)
	md[k] = []string{val}
}

// Join joins any number of mds into a single MD.
// The order of values for each key is determined by the order in which
// the mds containing those values are presented to Join.
func Join(mds ...MD) MD {
	out := MD{}
	for _, md := range mds {
		for k, v := range md {
			out[k] = append(out[k], v...)
		}
	}
	return out
}

// NewIncomingContext creates a new context with incoming md attached.
func NewIncomingContext(ctx context.Context, md MD) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.MD(md))
}

// NewOutgoingContext creates a new context with outgoing md attached. If used
// in conjunction with AppendToOutgoingContext, NewOutgoingContext will
// overwrite any previously-appended metadata.
func NewOutgoingContext(ctx context.Context, md MD) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.MD(md))
}

func FromIncomingContext(ctx context.Context) MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		return MD(md)
	}
	return nil
}

func FromOutgoingContext(ctx context.Context) MD {
	md, _, ok := metadata.FromOutgoingContextRaw(ctx)
	if ok {
		return MD(md)
	}

	return nil
}
