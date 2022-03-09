package main

import (
	. "file-server/utils"
	"net"
)

type Server struct {
	Network  string
	Address  string
	listener net.Listener
	Clients  []Client
	Channels []Channel
}

func NewServer() *Server {
	return &Server{
		Network: "tcp",
		Address: ":9999",
	}
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", ":9999")
	server.listener = listener
	if err != nil {
		PrintError(err.Error())
		return
	}

	server.Channels = append(server.Channels, *NewChannel("job"))
	server.Channels = append(server.Channels, *NewChannel("friends"))
	server.Channels = append(server.Channels, *NewChannel("homework"))

	go server.listen()
	PrintSuccess("----Server started----")
}

func (server *Server) listen() {
	disconnected := make(chan *Client)
	go server.disconnectClient(disconnected)
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			PrintError(err.Error())
			return
		}
		client := NewClient(connection)
		for i := 0; i < len(server.Channels); i++ {
			client.Register(&server.Channels[i])
		}

		server.Clients = append(server.Clients, *client)
		go client.handle(disconnected)
	}
}

func (server *Server) disconnectClient(disconnected chan *Client) {
	for client := range disconnected {
		for i, cl := range server.Clients {
			if cl.ID == client.ID {
				server.Clients = append(server.Clients[:i], server.Clients[i+1:]...)
				PrintWarning("Removed client")
			}
		}
		for _, channel := range server.Channels {
			channel.OnDisconnect(client)
		}

	}
}

func (server *Server) OnEntry(options []string) {
	switch options[0] {
	case START:
		server.Start()
	case STOP:
		server.Stop()
	}
}
func (server *Server) Identifier() string {
	return "server"
}

func (server *Server) Stop() {
	for i := 0; i < len(server.Clients); i++ {
		client := server.Clients[i]
		client.Connection.Close()
	}
	server.listener.Close()
}
