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
		apiId= 901789
		apiHash= "325ee9de98e6e16d3a615cb9331774e4"
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
			//fmt.Printf("%s\n", text.Text.Text)

			if Event.Message.SenderUserId != 780484786 {
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
				fmt.Printf("\nSendMessage response : %v \n %T \n", resp, resp)
				//resp.SendingState.MessageSendingStateType()
				//		}
			}
		}
	}
}

//149.154.167.40:443
//149.154.167.50:443  id: 552563440