package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/clava1096/rocket-service/notification/internal/config"
	"github.com/clava1096/rocket-service/notification/internal/service"
	"github.com/clava1096/rocket-service/notification/internal/service/consumer/order_assembled_consumer"
	"github.com/clava1096/rocket-service/notification/internal/service/consumer/order_paid_consumer"
	"github.com/clava1096/rocket-service/notification/internal/service/telegram"
	"github.com/clava1096/rocket-service/platform/pkg/closer"
	wrappedKafkaConsumer "github.com/clava1096/rocket-service/platform/pkg/kafka/consumer"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	kafkaMiddleware "github.com/clava1096/rocket-service/platform/pkg/middleware/kafka"
	"github.com/go-telegram/bot"

	httpClient "github.com/clava1096/rocket-service/notification/internal/client/http"
	telegramClient "github.com/clava1096/rocket-service/notification/internal/client/http/telegram"
	kafkaConverter "github.com/clava1096/rocket-service/notification/internal/converter/kafka"
	kafkaWrappedConverter "github.com/clava1096/rocket-service/notification/internal/converter/kafka/decoder"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
)

type diContainer struct {
	telegramService service.TelegramService
	telegramBot     *bot.Bot
	telegramClient  httpClient.TelegramClient

	consumerPaidGroup      sarama.ConsumerGroup
	consumerAssembledGroup sarama.ConsumerGroup
	orderAssembledConsumer kafka.Consumer
	orderPaidConsumer      kafka.Consumer

	orderAssemblyConsumerService service.OrderAssemblyConsumer
	orderPaidConsumerService     service.OrderPaidConsumer

	orderAssembledDecoder kafkaConverter.OrderAssemblyDecoder
	orderPaidDecoder      kafkaConverter.OrderPaidDecoder
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderPaidConsumer() service.OrderPaidConsumer {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = order_paid_consumer.NewService(
			d.KafkaOrderPaidConsumerInstance(),
			d.OrderPaidDecoder(),
			d.TelegramService(),
		)
	}
	return d.orderPaidConsumerService
}

func (d *diContainer) OrderAssembledConsumer() service.OrderAssemblyConsumer {
	if d.orderAssemblyConsumerService == nil {
		d.orderAssemblyConsumerService = order_assembled_consumer.NewService(
			d.KafkaOrderAssembledConsumerInstance(),
			d.OrderAssembledDecoder(),
			d.TelegramService(),
		)
	}
	return d.orderAssemblyConsumerService
}

func (d *diContainer) TelegramService() service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegram.NewService(
			d.TelegramClient(),
		)
	}
	return d.telegramService
}

func (d *diContainer) TelegramClient() httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot())
	}
	return d.telegramClient
}

func (d *diContainer) TelegramBot() *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().TelegramConfig.Token())
		if err != nil {
			panic(fmt.Sprintf("error creating telegram bot: %v", err))
		}

		d.telegramBot = b

	}
	return d.telegramBot
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = kafkaWrappedConverter.NewOrderPaidDecoder()
	}
	return d.orderPaidDecoder
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.OrderAssemblyDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = kafkaWrappedConverter.NewShipAssembledDecoder()
	}
	return d.orderAssembledDecoder
}

func (d *diContainer) KafkaOrderPaidConsumerInstance() kafka.Consumer {

	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerPaidGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumerConfig.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()))
	}

	return d.orderPaidConsumer
}

func (d *diContainer) ConsumerPaidGroup() sarama.ConsumerGroup {
	if d.consumerPaidGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().KafkaConfig.Brokers(),
			config.AppConfig().OrderPaidConsumerConfig.GroupID(),
			config.AppConfig().OrderPaidConsumerConfig.Config())

		if err != nil {
			panic(fmt.Sprintf("error creating consumer paid group: %v", err))
		}
		closer.AddNamed("Kafka order paid consumer group", func(ctx context.Context) error { return d.consumerPaidGroup.Close() })

		d.consumerPaidGroup = consumerGroup
	}
	return d.consumerPaidGroup
}

func (d *diContainer) KafkaOrderAssembledConsumerInstance() kafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerAssembledGroup(),
			[]string{
				config.AppConfig().OrderAssemblyConsumerConfig.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()))
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) ConsumerAssembledGroup() sarama.ConsumerGroup {
	if d.consumerAssembledGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().KafkaConfig.Brokers(),
			config.AppConfig().OrderPaidConsumerConfig.GroupID(),
			config.AppConfig().OrderAssemblyConsumerConfig.Config())

		if err != nil {
			panic(fmt.Sprintf("error creating consumer assembly group: %v", err))
		}

		d.consumerAssembledGroup = consumerGroup
	}
	return d.consumerAssembledGroup
}
