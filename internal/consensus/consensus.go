package consensus

import (
	"io"
	"io/ioutil"
	"sync"
	"github.com/hashicorp/raft"
)

const (
	userMaxByte = 64
	contentMaxByte = 1024
)

type Message struct {
	user    string
	content string
}

type MessageTracker struct {
	mtx  	 sync.RWMutex
	messages []*Message
}

type Snapshot struct {
	messages []*Message
}

var _ raft.FSM := &MessageTracker{}

func (m *Message) Init(message []byte) {
	m.user = string(message[0:userMaxByte])
	m.content = string(message[userMaxByte:])
}

func copyMessage(m *Message) *Message {
	copy := new(Message)
	copy.user = m.user
	copy.content = m.content
	return copy
}

func cloneMessages(messages []*Message) []*Message {
	var clone[]*Message
	for idx, m := range messages {
		clone = append(clone, copyMessage(m))
	}
	return clone
}

func concatenateMessages(messages []*Message) string {
	var concat string
	for _, message := range messages {
		user := message.user
		content := message.content
		filledUser = fillString(user, userMaxByte)
		filledContent := fillString(content, contentMaxByte)
		concat = append(concat, filledUser + filledContent)
	}
	return concat
}

func fillString(needFill string, maxByte int) string {
	escape := "\000"
	length := len(needFill)
	fill := string(needFill)
	for i := 0; i < (maxByte / 2) - length; i++ {
		fill += escape
	}
	return fill
}

func (mt *MessageTracker) Apply(l *raft.Log) interface{} {
	mt.mtx.Lock()
	defer mt.mtx.Unlock()
	message := new(Message)
	message.Init(l.Data) // Is l.Data a row of snapshot sink?
	mt.messages = append(mt.messages, message)	
	return nil
}

func (mt *MessageTracker) Snapshot() (raft.FSMSnapshot, error) {
	s := new(Snapshot)
	s.messages = cloneMessages(mt.messages)
	return &s, nil
}

func (mt *MessageTracker) Restore(r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	var messages[]*Message
	numEntries := len(b) / (userMaxByte + contentMaxByte)
	for i := 0; i < numEntries; i++ {
		offset := contentMaxByte * i
		user := b[offset:offset + userMaxByte]
		content := b[offset + userMaxByte:offset + contentMaxByte]
		message := new(Message)
		message.user = user
		message.content = content
		messages = append(messages, message)
	}
	mt.mtx.Lock()
	defer mt.mtx.Unlock()
	mt.messages = cloneMessages(messages)
	return nil
}

func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	_, err := sink.Write([]byte(concatenateMessages(s.messages)))
	if err != nil {
		sink.Cancel()
		return fmt.Errorf("sink.Write(): %v", err)
	}
	return sink.Close()
}

func (s *Snapshot) Release() { }