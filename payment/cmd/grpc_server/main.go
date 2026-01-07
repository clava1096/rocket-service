package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "github.com/clava1096/rocket-service/payment/internal/api/payment/v1"
	logger "github.com/clava1096/rocket-service/payment/internal/middleware"
	paymentService "github.com/clava1096/rocket-service/payment/internal/service/payment"
	paymentv1 "github.com/clava1096/rocket-service/shared/pkg/proto/payment/v1"
)

const grpcPort = ":50052"

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v", cerr)
		}
	}()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(logger.LoggerInterceptor()))

	service := paymentService.NewService()
	api := v1.NewAPI(service)
	paymentv1.RegisterPaymentServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("Starting gRPC server on port %s", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped")
}
