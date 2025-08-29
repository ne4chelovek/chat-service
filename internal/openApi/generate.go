package openApi

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i github.com/ne4chelovek/chat_service/internal/openApi.ApiCat -o ./mocks/ -s "_minimock.go"
