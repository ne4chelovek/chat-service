package chat

import (
	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"
	"strconv"
)

func (s *Server) ConnectChat(req *desc.ConnectChatRequest, stream desc.Chat_ConnectChatServer) error {
	id, err := strconv.Atoi(req.ChatId)
	if err != nil {
		return err
	}
	err = s.chatService.Connect(int64(id), req.GetUsername(), stream)
	if err != nil {
		return err
	}

	return err
}
