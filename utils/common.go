package utils

import (
	"context"
	"fmt"
	"github.com/aka-yz/go-micro-core/configs/log"
	"runtime"
)

func GoWithRecover(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			PrintStack(context.Background(), err)
		}
	}()

	fn()
}
func PrintStack(ctx context.Context, err interface{}) {
	//打印调用栈信息
	buf := make([]byte, 8192)
	n := runtime.Stack(buf, false)
	stackInfo := fmt.Sprintf("%s", buf[:n])
	log.Errorf(ctx, "err: %+v, err: panic stack info %s", err, stackInfo)
}
