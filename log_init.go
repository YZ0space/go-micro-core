package go_micro_core

import (
	"fmt"
	log2 "github.com/aka-yz/go-micro-core/configs/log"
	"github.com/aka-yz/go-micro-core/providers/constants"
	"go.uber.org/config"
)

func initLog(conf config.Provider) {
	var cv config.Value
	if cv = conf.Get(constants.ConfigKeyLog); !cv.HasValue() {
		return
	}

	var cfg log2.Option
	if err := cv.Populate(&cfg); err != nil {
		panic(err)
	}
	fmt.Printf("cfg:%v monitor\n", cfg)
	log2.InitLogger(&cfg)
}
