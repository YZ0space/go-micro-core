package extension

import "github.com/gin-gonic/gin"

type GinHandler interface {
	HandlerList() []*GinHandlerRegister
	MiddlewareList() []gin.HandlerFunc
}

type GinHandlerRegister struct {
	HttpMethod   string
	RelativePath string
	Handlers     []gin.HandlerFunc
}
