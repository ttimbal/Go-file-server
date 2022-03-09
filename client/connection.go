package main

import (
	. "file-server/utils"
	"net"
)

type Connection struct {
	Network  string
	Address  string
	Conn     net.Conn
	emitter  *Emitter
	listener *Listener
}

func NewConnection() *Connection {
	return &Connection{
		Network: "tcp",
		Address: ":9999",
	}
}

func (connection *Connection) Start() error {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		PrintError(err.Error())
		return err
	}
	PrintSuccess("----Connected to server----")
	connection.Conn = c
	connection.emitter = NewEmitter(connection.Conn)
	connection.listener = NewListener(connection.Conn)
	go connection.emitter.subscriptionListener(connection.listener.Responses)
	return nil
}

func (connection *Connection) HandleResponse() {
	connection.listener.Listen()
}

func (connection *Connection) close() {
	connection.listener.Stop()
	connection.Conn.Close()
}
