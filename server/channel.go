package main

import (
	"encoding/json"
	. "file-server/structs"
	. "file-server/utils"
)

type Observer interface {
	OnMessage(client *Client, message Message)
	OnDisconnect(client *Client)
	Identifier() string
}

type Channel struct {
	Name        string
	Subscribers []*Client
	broadcaster *Broadcaster
}

func NewChannel(name string) *Channel {
	return &Channel{
		Name:        name,
		broadcaster: BroadcastInstance(),
	}
}

func (channel *Channel) Subscribe(client *Client) {
	position := -1
	for i, subscriber := range channel.Subscribers {
		if subscriber.ID == client.ID {
			position = i
			break
		}
	}

	if position != -1 {
		return
	}
	channel.Subscribers = append(channel.Subscribers, client)
	PrintSuccess("New subscriber to", channel.Name, "Channel")

	msg, _ := json.Marshal("Successfully subscribed to " + channel.Name + " Channel")
	message := Message{Action: SUBSCRIBE, Channel: channel.Name, Message: msg}
	channel.broadcaster.Broadcast([]*Client{client}, message)
}

func (channel *Channel) Unsubscribe(client *Client) {
	position := -1
	for i, subscriber := range channel.Subscribers {
		if subscriber.ID == client.ID {
			position = i
			break
		}
	}

	if position == -1 {
		return
	}

	channel.Subscribers = append(channel.Subscribers[:position], channel.Subscribers[position+1:]...)

	msg, _ := json.Marshal("Successfully unsubscribed to " + channel.Name + " Channel")
	message := Message{Action: UNSUBSCRIBE, Channel: channel.Name, Message: msg}
	channel.broadcaster.Broadcast([]*Client{client}, message)
	PrintWarning("1 subscriber unsubscribed to", channel.Name, "channel")

}

func (channel *Channel) Send(client *Client, message Message) {
	position := -1
	for i, subscriber := range channel.Subscribers {
		if subscriber.ID == client.ID {
			position = i
			break
		}
	}

	if position == -1 {
		return
	}

	var fileMessage FileMessage
	err := json.Unmarshal(message.Message, &fileMessage)
	if err != nil {
		PrintError(err.Error())
		return
	}
	PrintSuccess("A received file to", channel.Name, "Channel")

	_copy := make([]*Client, len(channel.Subscribers))
	copy(_copy, channel.Subscribers)
	_copy = append(_copy[:position], _copy[position+1:]...)
	channel.broadcaster.Broadcast(_copy, message)

	msg, _ := json.Marshal("Successfully sent to " + channel.Name + " Channel")
	newMessage := Message{Action: SEND, Channel: channel.Name, Message: msg}
	channel.broadcaster.Broadcast([]*Client{client}, newMessage)

}

func (channel *Channel) OnMessage(client *Client, message Message) {
	switch message.Action {
	case SUBSCRIBE:
		if channel.Name == message.Channel {
			go channel.Subscribe(client)
		}
	case UNSUBSCRIBE:
		go channel.Unsubscribe(client)
	case SEND:
		go channel.Send(client, message)
	}
}

func (channel *Channel) OnDisconnect(client *Client) {
	channel.Unsubscribe(client)
}

func (channel *Channel) Identifier() string {
	return channel.Name
}
