package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/clava1096/rocket-service/assembly/internal/config"
	kafkaConverter "github.com/clava1096/rocket-service/assembly/internal/converter/kafka"
	"github.com/clava1096/rocket-service/assembly/internal/converter/kafka/decoder"
	"github.com/clava1096/rocket-service/assembly/internal/service"
	"github.com/clava1096/rocket-service/assembly/internal/service/consumer/order_consumer"
	"github.com/clava1096/rocket-service/assembly/internal/service/produser/order_producer"
	"github.com/clava1096/rocket-service/platform/pkg/closer"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/clava1096/rocket-service/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/clava1096/rocket-service/platform/pkg/kafka/producer"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
)

type diContainer struct {
	consumerGroup         sarama.ConsumerGroup
	orderAssemblyConsumer kafka.Consumer

	orderConsumerService service.OrderConsumerService
	orderProducerService service.OrderProducerService
	orderPaidDecoder     kafkaConverter.OrderPaidDecoder
	orderProducer        kafka.Producer
	syncProducer         sarama.SyncProducer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderConsumerService() service.OrderConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = order_consumer.NewService(d.OrderAssemblyConsumer(), d.orderProducerShipAssembled(), d.OrderPaidDecoder())
	}
	return d.orderConsumerService
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) OrderAssemblyConsumer() kafka.Consumer {
	if d.orderAssemblyConsumer == nil {
		d.orderAssemblyConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger())
	}
	return d.orderAssemblyConsumer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config())
		if err != nil {
			panic(fmt.Sprintf("Failed to create consumer group: %s\n", err.Error()))
		}
		d.consumerGroup = consumerGroup
	}
	return d.consumerGroup
}

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = order_producer.NewService(d.orderProducerShipAssembled())
	}
	return d.orderProducerService
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledProducer.Config())
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error { return p.Close() })
		d.syncProducer = p
	}
	return d.syncProducer
}

func (d *diContainer) orderProducerShipAssembled() kafka.Producer {
	if d.orderProducer == nil {
		d.orderProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledProducer.Topic(),
			logger.Logger())
	}
	return d.orderProducer
}
