package go_micro_core

import (
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	"go.uber.org/config"
)

var conf config.Provider

var (
	configPath = flag.String("p", "./configs/application.yml", "configs for service")
	prefix     string
)

func SetConfigPathPrefix(pathPrefix string) {
	prefix = pathPrefix
}

func initConfig() {
	path := *configPath
	if !Env.Dev() {
		path = fmt.Sprintf(".%v/configs/application-%v.yml", prefix, Env)
	}

	var err error
	if conf, err = config.NewYAML(
		config.File(path),
	); err != nil {
		panic(err)
	}
	return
}

func LoadAppConf(key string, c interface{}) {
	if err := conf.Get(key).Populate(c); err != nil {
		panic(err)
	}
	y, _ := yaml.Marshal(c)
	fmt.Printf("[Config] LoadConf \n%+s\n", y)
}
