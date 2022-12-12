# go-micro-core


name = github.com/aka-yz/go-micro-core

version = go1.18

micro service core module, include:
- object injecting based facebookgo/inject.go
- object life cycle management(init -> starter -> stop)
- local env-param read
- app config(yaml, based viper) read & config
- mysql, pg, redis connection initial


use it:

````go get -u "github.com/aka-yz/go-micro-core@v0.0.4"````

redis-client: https://github.com/go-redis/redis/v8
postgresql-client: https://github.com/go-pg/pg
mysql-client: https://github.com/gocraft/dbr/v2
log: https://github.com/rs/zerolog
