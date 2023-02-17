package network

import (
	"bytes"
	"encoding/json"
)

const (
	GetChain = iota
	GetBlocks
	MsgName3
	MsgName4
)

type Message struct {
	Id       int
	StreamId string
	payload  []byte
}

func NewMessage(id int, payload []byte) *Message {
	return &Message{
		Id:      id,
		payload: payload,
	}
}

func Marhal(msg *Message) ([]byte, error) {
	return json.Marshal(struct {
		Id       int
		StreamId string
		Payload  []byte
	}{
		Id:       msg.Id,
		StreamId: msg.StreamId,
		Payload:  msg.payload,
	})
}

func Unmarhsal(d []byte) (*Message, error) {
	unmarshalled := struct {
		Id      int
		Payload []byte
	}{}
	err := json.Unmarshal(d, &unmarshalled)
	if err != nil {
		return nil, err
	}
	return &Message{
		Id:      unmarshalled.Id,
		payload: unmarshalled.Payload,
	}, nil
}

func NilMessage() *Message {
	return &Message{
		payload: []byte("nil"),
	}
}

func IsNilMessage(mes *Message) bool {
	if bytes.Equal(mes.payload, []byte("nil")) {
		return true
	}
	return false
}

func (msg *Message) Payload() []byte {
	return msg.payload
}
