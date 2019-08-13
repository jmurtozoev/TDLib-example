package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/zelenin/go-tdlib/client"
)

func main() {
	// client authorizer
	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)

	const (
		apiId= 00000 // telegram api ID
		apiHash= "ksdhfjkhsjfhfkjnsdkcjfnsdkfkdjsnf" // telegram api hash key 
	)

	authorizer.TdlibParameters <- &client.TdlibParameters{
		UseTestDc:              false,
		DatabaseDirectory:      filepath.Join(".tdlib", "database"),
		FilesDirectory:         filepath.Join(".tdlib", "files"),
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  apiId,
		ApiHash:                apiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Server",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}
	//Receive updates
	tdlibClient, err := client.NewClient(authorizer)
	if err != nil {
		log.Fatalf("NewClient error: %s", err)
	}

	listener := tdlibClient.GetListener()
	defer listener.Close()

	for update := range listener.Updates {
		if update.GetClass() == client.ClassUpdate && update.GetType() == client.TypeUpdateNewMessage {
			Event := update.(*client.UpdateNewMessage)
			text, _ := Event.Message.Content.(*client.MessageText)
			
			if Event.Message.SenderUserId != 000000000 {  // your account id
				req := client.SendMessageRequest{
					ChatId: Event.Message.ChatId,
					InputMessageContent: &client.InputMessageText{
						Text: &client.FormattedText{
							Text:     fmt.Sprintf("%s\nNumber of chars: %d", text.Text.Text, len(text.Text.Text)),
							Entities: nil,
						},
						ClearDraft: true,
					},
				}
				resp, err := tdlibClient.SendMessage(&req)
				if err != nil {
					log.Fatalf("SendMessage error: %s, \t%v\n", err, resp)
				}
			}
		}
	}
}
