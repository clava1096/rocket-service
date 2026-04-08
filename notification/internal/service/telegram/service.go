package telegram

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"time"

	"github.com/clava1096/rocket-service/notification/internal/client/http"
	"github.com/clava1096/rocket-service/notification/internal/model"
)

//go:embed templates/paid_notification.tmpl
var templatesPaid embed.FS

//go:embed templates/assembled_notification.tmpl
var templatesAssembled embed.FS

const chatID = 123 //todo придумать механизм чтобы чат айди сохранялись

type paidTemplateData struct {
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
	PaidAt          time.Time
}

type assembledTemplateData struct {
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
	AssembledAt  time.Time
}

var orderPaidTemplate = template.Must(template.ParseFS(templatesPaid, "templates/paid_notification.tmpl"))
var orderAssemblyTemplate = template.Must(template.ParseFS(templatesAssembled, "templates/assembled_notification.tmpl"))

type service struct {
	telegramClient http.TelegramClient
}

func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) SendPaidNotification(ctx context.Context, event model.OrderPaidEvent) error {
	message, err := s.buildPaidMessage(event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) buildPaidMessage(event model.OrderPaidEvent) (string, error) {
	data := paidTemplateData{
		OrderUUID:       event.OrderUUID,
		UserUUID:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUUID: event.TransactionUUID,
		PaidAt:          time.Now(),
	}

	var buf bytes.Buffer
	err := orderPaidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) SendAssembledNotification(ctx context.Context, event model.ShipAssembledEvent) error {
	message, err := s.buildAssembledMessage(event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) buildAssembledMessage(event model.ShipAssembledEvent) (string, error) {
	data := assembledTemplateData{
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
		AssembledAt:  time.Now(),
	}

	var buf bytes.Buffer
	err := orderAssemblyTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
