package interceptors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aka-yz/go-micro-core/configs/log"
	"github.com/aka-yz/go-micro-core/providers/transport/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"runtime"
	"strings"
	"time"
)

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		now := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			log.Infof(ctx, "method:%s req:%v reply:%v err:%v elapsed:%v", getMethod(method), req, logCutOff(reply), err, time.Since(now))
		} else {
			log.Errorf(ctx, "method:%s req:%v reply:%v err:%v elapsed:%v", getMethod(method), req, logCutOff(reply), err, time.Since(now))
		}

		return err
	}
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if pErr := recover(); pErr != nil {
				printStack(ctx, req, pErr)
				err = errors.New("grpc server panic")
			}
		}()

		now := time.Now()
		resp, err = handler(ctx, req)
		if err == nil {
			log.Infof(ctx, "method:%s req:%v reply:%v err:%v elapsed:%v", getMethod(info.FullMethod), req, resp, err, time.Since(now))
		} else {
			log.Errorf(ctx, "method:%s req:%v reply:%v err:%v elapsed:%v", getMethod(info.FullMethod), req, resp, err, time.Since(now))
		}

		return resp, err
	}
}

func GinUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := gin.MetadataFromContext(ctx)
		if md != nil {
			ctx = metadata.NewOutgoingContext(ctx, metadata.MD(md))
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func GinUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ctx = gin.NewContextFromMetadata(ctx, gin.Metadata(md))
		}
		return handler(ctx, req)
	}
}

func getMethod(method string) string {
	ms := strings.Split(method, "/")
	return ms[len(ms)-1]
}

const (
	LogLenLimit = 3072 // 超过4096会打到daemon
)

func logCutOff(reply interface{}) interface{} {
	replyTmp, err := json.Marshal(reply)
	if err != nil {
		return reply
	}
	replyLog := string(replyTmp)
	if len(replyLog) > LogLenLimit {
		replyLog = replyLog[:LogLenLimit] + "..."
	}
	return replyLog
}

func printStack(ctx context.Context, req interface{}, err interface{}) {
	//打印调用栈信息
	buf := make([]byte, 8192)
	n := runtime.Stack(buf, false)
	stackInfo := fmt.Sprintf("%s", buf[:n])
	log.Errorf(ctx, "err: %v, req: %+v panic stack info %s", err, req, stackInfo)
}
