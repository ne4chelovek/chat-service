package converter

import (
	model "github.com/ne4chelovek/chat_service/internal/model"
	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"
)

func ToChatFromDesc(chat *desc.Message) *model.Message {
	return &model.Message{
		From:        chat.From,
		Text:        chat.Text,
		Timestamppb: chat.CratedAt,
	}
}

func ToMessageFromService(message *model.Message) *desc.Message {
	return &desc.Message{
		From:     message.From,
		Text:     message.Text,
		CratedAt: message.Timestamppb,
	}
}
