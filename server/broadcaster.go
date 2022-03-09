package main

import (
	"encoding/json"
	. "file-server/structs"
	. "file-server/utils"
)

type Broadcaster struct {
}

var singleInstance *Broadcaster

func BroadcastInstance() *Broadcaster {
	if singleInstance == nil {
		singleInstance = &Broadcaster{}
	}
	return singleInstance
}

func (broadcaster *Broadcaster) Broadcast(subscribers []*Client, message Message) {
	for _, subscriber := range subscribers {
		broadcaster.emit(subscriber, message)
	}
}

func (broadcaster *Broadcaster) emit(subscriber *Client, message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		PrintError(err.Error())
	}
	if len(data) > MAX_SIZE {
		PrintWarning("You can not upload more than 5 MB")
		return
	}
	subscriber.Connection.Write(data)
}
