package middleware

import (
	"crypto/tls"
	"github.com/aka-yz/go-micro-core/providers/option"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewClient(opt *option.Redis) redis.UniversalClient {
	var redisClient redis.UniversalClient
	if opt.IsClusterMode {
		redisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        opt.ClusterAddr,
			MaxRedirects: opt.MaxRetries,
			ReadTimeout:  time.Duration(opt.ReadTimeout),
			PoolSize:     opt.PoolSize,
			Password:     opt.Password,
			TLSConfig: &tls.Config{
				// nolint: gosec
				InsecureSkipVerify: true,
			},
		})
	} else if opt.IsElastiCache {
		redisClient = redis.NewClient(&redis.Options{
			Addr:        opt.Addr,
			ReadTimeout: time.Duration(opt.ReadTimeout),
			PoolSize:    opt.PoolSize,
			Password:    opt.Password,
			TLSConfig: &tls.Config{
				// Elasticache cert cannot be applied to cname record we use
				// nolint: gosec
				InsecureSkipVerify: true,
			},
		})
	} else {
		redisClient = redis.NewClient(&redis.Options{
			Addr:        opt.Addr,
			ReadTimeout: time.Duration(opt.ReadTimeout),
			PoolSize:    opt.PoolSize,
			Password:    opt.Password,
		})
	}
	return redisClient
}
