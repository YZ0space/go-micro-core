package go_micro_core

import (
	"fmt"
	"github.com/aka-yz/go-micro-core/providers/config/log"
	"github.com/aka-yz/go-micro-core/providers/constants"
	"go.uber.org/config"
)

func initLog(conf config.Provider) {
	var cv config.Value
	if cv = conf.Get(constants.ConfigKeyLog); !cv.HasValue() {
		return
	}

	var cfg log.Option
	if err := cv.Populate(&cfg); err != nil {
		panic(err)
	}
	fmt.Printf("cfg:%v monitor\n", cfg)
	log.InitLogger(&cfg)
}
