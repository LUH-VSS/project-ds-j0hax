package lib

import "time"

type Message struct {
	Timestamp time.Time
	Words     []string
}

func New(words []string) *Message {
	return &Message{
		Timestamp: time.Now(),
		Words:     words,
	}
}
