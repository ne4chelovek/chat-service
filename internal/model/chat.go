package model

import "google.golang.org/protobuf/types/known/timestamppb"

type Message struct {
	From        string
	Text        string
	Timestamppb *timestamppb.Timestamp
}
