package model

import (
	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Message struct {
	From        string
	Text        string
	Timestamppb *timestamppb.Timestamp
}

type WsMessage struct {
	ChatID int64  `json:"chat_id"`
	From   string `json:"from"`
	Text   string `json:"text"`
}

type Stream interface {
	desc.Chat_ConnectChatServer
}
