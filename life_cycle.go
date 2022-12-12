package go_micro_core

// 当前服务生命周期抽象
type initialization interface {
	Init()
}

type starter interface {
	Start()
}

type stoper interface {
	Stop()
}
