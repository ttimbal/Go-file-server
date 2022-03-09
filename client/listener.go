package main

import (
	"encoding/json"
	. "file-server/structs"
	. "file-server/utils"
	"net"
	"os"
)

type Listener struct {
	Connection net.Conn
	Responses  map[string]chan string
	active     bool
}

func NewListener(connection net.Conn) *Listener {
	responses := make(map[string]chan string)
	responses[SUBSCRIBE] = make(chan string)
	responses[UNSUBSCRIBE] = make(chan string)
	return &Listener{
		Connection: connection,
		Responses:  responses,
		active:     true,
	}
}

func (listener *Listener) Listen() {
	var response Message
	for listener.active {
		b := make([]byte, MAX_SIZE)
		bs, err := listener.Connection.Read(b)
		if err != nil {
			PrintError(err.Error(), "Disconnected")
			listener.Stop()
			break
		} else {
			err = json.Unmarshal(b[:bs], &response)
			if err != nil {
				PrintError(err.Error())
				continue
			}

			switch response.Action {
			case SUBSCRIBE:
				listener.Subscribe(response)
			case UNSUBSCRIBE:
				listener.Unsubscribe(response)
			case SEND:
				listener.Send(response)
			}
		}
	}
}

func (listener *Listener) Stop() {
	listener.active = false
}

func (listener *Listener) Subscribe(response Message) {
	PrintSuccess(string(response.Message))
	listener.Responses[SUBSCRIBE] <- response.Channel
}

func (listener *Listener) Unsubscribe(response Message) {
	PrintSuccess(string(response.Message))
	listener.Responses[UNSUBSCRIBE] <- response.Channel
}

func (listener *Listener) Send(response Message) {
	var fileMessage FileMessage
	err := json.Unmarshal(response.Message, &fileMessage)
	if err != nil {
		//PrintError(err.Error())
		PrintSuccess(string(response.Message))
		return
	}
	err = os.WriteFile("../storage/"+fileMessage.Name, fileMessage.Content, fileMessage.Mode)
	PrintSuccess("New file saved in Storage")
	if err != nil {
		PrintError("Error saving file", err.Error())
	}

}
