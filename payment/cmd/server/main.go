package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	logger "github.com/clava1096/rocket-service/payment/internal"
	paymentv1 "github.com/clava1096/rocket-service/shared/pkg/proto/payment/v1"
)

const grpcPort = ":50052"

type paymentServer struct {
	paymentv1.UnimplementedPaymentServiceServer
}

func (s *paymentServer) PayOrder(ctx context.Context, payOrder *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	transactionUuid := uuid.New()
	log.Printf("Payment uuid: %s | order_uuid: %s | user_uuid: %s | payment_method: %s", transactionUuid,
		payOrder.OrderUuid,
		payOrder.UserUuid,
		payOrder.PaymentMethod.String())
	return &paymentv1.PayOrderResponse{TransactionUuid: transactionUuid.String()}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost"+grpcPort)
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

	service := &paymentServer{}
	paymentv1.RegisterPaymentServiceServer(s, service)

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
