FROM golang:1.25.0-alpine as builder
LABEL authors="ant1soc1al"

RUN apk add --no-cache git

WORKDIR /app

COPY go.work go.work.sum

COPY platform/go.mod platform/go.sum ./platform/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY shared/go.mod shared/go.sum ./shared/
COPY order/go.mod order/go.sum ./order/
COPY payment/go.mod payment/go.sum ./payment/

RUN cd order && go mod download

COPY platform ./platform
COPY inventory ./inventory
COPY shared ./shared
COPY order ./order
COPY payment ./payment

#ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 ./grpc-health-probe
#todo написать wget для этого сервиса вместо grpc_health.
#RUN chmod +x grpc-health-probe

RUN cd order && CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o ../app-order ./cmd/http_server/main.go


FROM alpine:3.21.3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

USER appuser

COPY --from=builder /app/app-order .

COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe

COPY deploy/compose/order/.env .env

EXPOSE 8080

ENTRYPOINT ["./app-order"]
