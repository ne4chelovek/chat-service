package converter

import (
	model "github.com/ne4chelovek/chat_service/internal/model"
	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"
)

func ToChatFromDesc(chat *desc.Message) *model.Message {
	return &model.Message{
		From: chat.From,
		Text: chat.Text,
	}
}

func ToMessageFromService(message *model.Message) *desc.StreamMessage {
	return &desc.StreamMessage{
		From:      message.From,
		Text:      message.Text,
		CreatedAt: message.Timestamppb,
	}
}
