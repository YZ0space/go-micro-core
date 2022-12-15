package extension

import "github.com/gin-gonic/gin"

// GinHandler 利用接口特性处理 gin-handler
type GinHandler interface {
	HandlerList() []*GinHandlerRegister
	MiddlewareList() []gin.HandlerFunc
}

type GinHandlerRegister struct {
	HttpMethod   string
	RelativePath string
	Handlers     []gin.HandlerFunc
}
