package chat

import (
	"context"
	"github.com/ne4chelovek/chat_service/internal/converter"
	"github.com/ne4chelovek/chat_service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

func (s *serv) Connect(chatID int64, username string, stream model.Stream) error {
	id := strconv.Itoa(int(chatID))
	s.mxChannels.RLock()
	chatChan, ok := s.channels[id]
	s.mxChannels.RUnlock()

	if !ok {
		return status.Errorf(codes.NotFound, "chat not found")
	}

	s.mxChats.Lock()
	if _, okChat := s.chats[id]; !okChat {
		s.chats[id] = &chat{
			streams: make(map[string]model.Stream),
		}
	}
	s.mxChats.Unlock()

	s.chats[id].m.Lock()
	s.chats[id].streams[username] = stream
	s.chats[id].m.Unlock()

	// если в чате есть сообщения, отправляем первые 10, если их нет придёт nil nil
	var page uint64 = 0
	err := s.sendHistory(stream.Context(), chatID, stream, page)
	if err != nil {
		return err
	}

	for {
		select {
		case msg, okCh := <-chatChan:
			if !okCh {
				return nil
			}
			if strings.HasPrefix(msg.Text, "!loadmore") {
				page++
				if err := s.sendHistory(stream.Context(), chatID, stream, page); err != nil {
					return err
				}
			}
			if strings.HasPrefix(msg.Text, "!cat") {
				fact, err := s.api.GetCatFact()
				if err != nil {
					return err
				}
				err = stream.Send(converter.ToMessageFromService(fact))
				if err != nil {
					return err
				}
			}

			err := s.sendingMessages(id, msg)
			if err != nil {
				return err
			}

		case <-stream.Context().Done():
			s.chats[id].m.Lock()
			delete(s.chats[id].streams, username)
			s.chats[id].m.Unlock()
			return nil
		}
	}
}

func (s *serv) sendHistory(ctx context.Context, chatId int64, stream model.Stream, page uint64) error {
	messages, err := s.chatRepository.GetMessage(ctx, chatId, page)
	if err != nil {
		return err
	}
	for _, msg := range messages {
		if err := stream.Send(converter.ToMessageFromService(msg)); err != nil {
			return err
		}
	}
	return nil
}

func (s *serv) sendingMessages(id string, msg *model.Message) error {
	s.mxChats.RLock()
	defer s.mxChats.RUnlock()

	chat, ok := s.chats[id]
	if !ok {
		return nil
	}

	for _, stream := range chat.streams {
		if err := stream.Send(converter.ToMessageFromService(msg)); err != nil {
			return err
		}
	}
	return nil
}
