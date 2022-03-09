package main

import (
	"encoding/json"
	. "file-server/structs"
	. "file-server/utils"
	"io/ioutil"
	"net"
	"os"
)

type Emitter struct {
	Connection    net.Conn
	subscriptions []string
}

func NewEmitter(connection net.Conn) *Emitter {
	return &Emitter{
		Connection: connection,
	}
}

func (emitter *Emitter) Subscribe(channel string) {
	if emitter.isSubscribed(channel) {
		PrintWarning("You are subscribed to", channel, "channel")
		return
	}
	message := Message{Action: SUBSCRIBE, Channel: channel}
	emitter.emit(message)
}

func (emitter *Emitter) Unsubscribe(channel string) {
	if !emitter.isSubscribed(channel) {
		PrintWarning("You are not subscribed to", channel, "channel")
		return
	}
	message := Message{Action: UNSUBSCRIBE, Channel: channel}
	emitter.emit(message)
}

func (emitter *Emitter) SendFile(channel string, filePath string) {
	if !emitter.isSubscribed(channel) {
		PrintWarning("You are not subscribed to", channel, "channel")
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		PrintError(err.Error())
		return
	}
	defer file.Close()
	infoFile, _ := file.Stat()
	fileByte, _ := ioutil.ReadAll(file)
	fileMessage := FileMessage{Name: infoFile.Name(), Size: infoFile.Size(), Content: fileByte, Mode: infoFile.Mode()}
	fileMessageByte, _ := json.Marshal(fileMessage)

	message := Message{Action: SEND, Channel: channel, Message: fileMessageByte}
	emitter.emit(message)
}

func (emitter *Emitter) emit(message Message) {
	data, _ := json.Marshal(message)
	if len(data) > MAX_SIZE {
		PrintWarning("You can not upload more than 5 MB")
		return
	}
	emitter.Connection.Write(data)
}

func (emitter *Emitter) isSubscribed(channel string) bool {
	isSubscribed := false
	for _, subscription := range emitter.subscriptions {
		if subscription == channel {
			isSubscribed = true
			break
		}
	}
	return isSubscribed
}

func (emitter *Emitter) OnEntry(options []string) {
	switch options[0] {
	case SUBSCRIBE:
		emitter.Subscribe(options[1])
	case UNSUBSCRIBE:
		emitter.Unsubscribe(options[1])
	case SEND:
		emitter.SendFile(options[1], options[2])
	case EXIT:
		break
	}

}
func (emitter *Emitter) Identifier() string {
	return "client"
}

func (emitter *Emitter) subscriptionListener(responses map[string]chan string) {
	go func() {
		for channel := range responses[SUBSCRIBE] {
			if !emitter.isSubscribed(channel) {
				emitter.subscriptions = append(emitter.subscriptions, channel)
			}
		}
	}()

	go func() {
		for response := range responses[UNSUBSCRIBE] {
			position := -1
			for i, subscription := range emitter.subscriptions {
				if subscription == response {
					position = i
					break
				}
			}
			if position != -1 {
				emitter.subscriptions = append(emitter.subscriptions[:position], emitter.subscriptions[position+1:]...)
			}
		}
	}()
}
