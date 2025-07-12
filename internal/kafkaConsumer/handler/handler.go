package handler

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

type MessageHandler interface {
	HandleMessage(message []byte, offset kafka.Offset) error
}
type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// HandleMessage просто логируем, но можно сохранять в бд
func (h *Handler) HandleMessage(message []byte, offset kafka.Offset) error {
	log.Printf(" Message from kafka whith offset: %d, %s", offset, string(message))
	return nil
}
