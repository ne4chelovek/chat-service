package apiHttp

import (
	"encoding/json"
	"github.com/ne4chelovek/chat_service/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"net/http"
)

const url = "https://catfact.ninja/fact"

type HttpClient struct {
	client *http.Client
}

func NewApiClient(client *http.Client) *HttpClient {
	return &HttpClient{
		client: client,
	}
}

func (c *HttpClient) GetCatFact() (*model.Message, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Неверный статус ответа: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
	}

	var fact *model.CatFact
	err = json.Unmarshal(body, &fact)
	if err != nil {
		log.Printf("Ошибка при парсинге JSON: %v", err)
	}

	mes := &model.Message{
		From:        "",
		Text:        fact.Fact,
		Timestamppb: timestamppb.Now(),
	}

	return mes, nil
}
