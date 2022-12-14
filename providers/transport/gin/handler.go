package gin

import (
	"context"
	"fmt"
	"github.com/aka-yz/go-micro-core/configs/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var reg, _ = regexp.Compile("/[0-9]+")

type Handler struct {
	*gin.Engine
	opts Options
}

func NewHandler(opts ...Option) *Handler {
	opt := Options{
		interceptors: []Interceptor{NopInterceptor},
		Codec:        defaultCodec{},
	}

	for _, o := range opts {
		o(&opt)
	}

	handler := &Handler{
		Engine: newEngine(),
		opts:   opt,
	}

	return handler
}

func newEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(Logger(), gin.Recovery())
	return engine
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}
		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path
		log.Infof(context.Background(), "[GIN] %3d| %13v | %15s |%-7s %#v |%s",
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage)

	}
}

func (s *Handler) RegisterService(srv interface{}) {
	// 注册当前服务
	service, err := register(srv)
	if err != nil {
		panic(err)
	}

	pe, err := newProtoExtend(s.opts.ProtoName)
	if err != nil {
		panic(err)
	}

	s.opts.ServiceName = service.name
	prefix := s.opts.Prefix
	if prefix == "" {
		prefix = fmt.Sprintf("/%s/", s.opts.ServiceName)
	}
	for method, v := range service.methods {
		rules := pe.methodHttpRules(method)
		for _, rule := range rules {
			s.Handle(rule.method, path.Join(prefix, rule.pattern), s.ginHandler(v))
		}
	}
}

func (s *Handler) Option() Options {
	return s.opts
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idx := strings.LastIndex(r.URL.Path, "/")
	if idx < 0 {
		panic(fmt.Errorf("rpc: no method in path %q", r.URL.Path))
	}

	if s.opts.Prefix == "" {
		r.URL.Path = "/" + s.opts.ServiceName + r.URL.Path[idx:]
	}

	s.Engine.ServeHTTP(w, r)
}

func (s *Handler) ginHandler(methodSpec *ServiceMethod) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		md := NewMetadata(r.Header, "sid")
		md.Set("request-time", fmt.Sprintf("%d", time.Now().UnixNano()))
		md.Set("method-name", methodSpec.method.Name)
		md.Set("client-ip", ctx.ClientIP())
		md = Join(md, Metadata(r.URL.Query()), ginParams(ctx.Params))

		newCtx := NewContextFromMetadata(r.Context(), md)
		r = r.WithContext(newCtx)

		reply, err := s.doHandle(r, methodSpec)
		if err != nil {
			s.opts.Codec.Encode(r, ctx.Writer, err)
		} else {
			s.opts.Codec.Encode(r, ctx.Writer, reply)
		}
		return
	}
}

func (s *Handler) doHandle(r *http.Request, methodSpec *ServiceMethod) (reply interface{}, err error) {
	ctx := r.Context()

	reqValue, err := s.getRequestValue(methodSpec.ReqType, r)
	if err != nil {
		return
	}

	return ChainInterceptors(s.opts.interceptors...)(ctx, reqValue, methodSpec.call)
}

func (s *Handler) getRequestValue(p reflect.Type, r *http.Request) (value interface{}, err error) {
	reqValue := reflect.New(p)
	err = s.opts.Codec.Decode(r, reqValue.Interface())
	return reqValue.Interface(), err
}

func ToUnderLine(src string) string {
	var dest []byte
	for _, s := range src {
		if s >= 'A' && s <= 'Z' {
			dest = append(dest, byte('_'), byte(s+32))
		} else {
			dest = append(dest, byte(s))
		}
	}

	if dest[0] == '_' {
		dest = dest[1:]
	}
	return string(dest)
}

func getMethodName(path string) (method string) {
	return path[strings.LastIndex(path, "/")+1:]
}

func ginParams(params gin.Params) Metadata {
	md := Metadata{}
	for _, param := range params {
		md.Set(param.Key, param.Value)
	}

	return md
}

func getSafeRequestURI(c *gin.Context) string {
	return c.Request.Method + "_" + reg.ReplaceAllString(c.Request.URL.Path, "/:Id")
}
