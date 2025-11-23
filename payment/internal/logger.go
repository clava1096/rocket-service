package internal

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	"path"
	"time"
)

func LoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		method := path.Base(info.FullMethod)

		log.Printf("Started gRPC method: %s", method)
		start := time.Now()

		resp, err = handler(ctx, req)
		duration := time.Since(start)

		// Форматируем сообщение в зависимости от результата
		if err != nil {
			st, _ := status.FromError(err)
			log.Printf("❌ Finished gRPC method %s with code %s: %v (took: %v)\n", method, st.Code(), err, duration)
		} else {
			log.Printf("✅ Finished gRPC method %s successfully (took: %v)\n", method, duration)
		}

		return resp, err
	}
}
