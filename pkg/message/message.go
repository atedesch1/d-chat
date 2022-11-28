package message

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Message struct {
	From int
	To   int
}

func NewMessage(fromId int, toId int) *Message {
	return &Message{
		From: fromId,
		To:   toId,
	}
}

func (m *Message) EncodeToBytes() []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func DecodeToMessage(s []byte) *Message {
	p := &Message{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
