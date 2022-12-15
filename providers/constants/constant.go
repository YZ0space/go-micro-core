package constants

const (
	ConfigKeyPostgresql = "pg"
	ConfigKeyMysql      = "db"
	ConfigKeyRedis      = "redis"
	ConfigKeyLog        = "log"
	ConfigHttpClient    = "httpclient"
	ConfigSrvKey        = "gin-server"
	ConfigApp           = "app"
	ConfigCron          = "cron-pool"
	HandlerInjectName   = "gin-handler"
	GinMethodGet        = "GET"
	GinMethodPost       = "POST"
)

// AllowedOrigins used in local test only
var AllowedOrigins = []string{
	"http://localhost:8080",
	"https://space.id",
	"https://app.space.id",
	"https://pre.stg.space.id",
	"https://app.stg.space.id",
	"https://alpha.space.id",
	"http://localhost:3000",
}
var AllowedHeaders = []string{
	"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "Signature",
}
