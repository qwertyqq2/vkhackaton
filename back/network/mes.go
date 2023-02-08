package network

import "encoding/json"

const (
	MsgName1 = iota
	MsgName2
	MsgName3
	MsgName4
)

type Message struct {
	Id      int
	payload []byte
}

func NewMessage(id int, payload []byte) *Message {
	return &Message{
		Id:      id,
		payload: payload,
	}
}

func Marhal(msg *Message) ([]byte, error) {
	return json.Marshal(struct {
		Id      int
		Payload []byte
	}{
		Id:      msg.Id,
		Payload: msg.payload,
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
