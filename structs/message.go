package structs

import (
	"encoding/json"
	"os"
)

type Message struct {
	Action  string          `json:"action"`
	Channel string          `json:"channel"`
	Message json.RawMessage `json:"message"`
}

type FileMessage struct {
	Name    string
	Size    int64
	Content []byte
	Mode    os.FileMode
}
