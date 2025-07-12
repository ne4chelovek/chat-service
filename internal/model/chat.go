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

type Stream interface {
	desc.Chat_ConnectChatServer
}
