package openApi

import (
	"github.com/ne4chelovek/chat_service/internal/model"
)

type ApiCat interface {
	GetCatFact() (*model.Message, error)
}
