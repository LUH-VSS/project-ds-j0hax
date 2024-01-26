package endpoint

import (
	"encoding/gob"
	"log"
	"net"
	"time"

	"github.com/LUH-VSS/project-ds-j0hax/lib"
)

type Endpoint struct {
	conn      net.Conn
	enc       *gob.Encoder
	wordQueue chan string
	done      chan bool
}

// sendBatch sends a collectiong of words to the configured endpoint.
//
// It will attempt to retransmit in case of an error.
func (e *Endpoint) sendBatch(batch []string) {
	backoff := time.Second
	for {
		data := lib.New(batch)
		err := e.enc.Encode(data)
		if err != nil {
			log.Print(err)
			log.Printf("Retrying in %s\n", backoff)
			time.Sleep(backoff)
			backoff += time.Second
		} else {
			return
		}
	}
}

// Run begins adding queued words to a batch
func (e *Endpoint) Run() {
	batch := make([]string, 0, cap(e.wordQueue))
	for word := range e.wordQueue {
		batch = append(batch, word)
		if len(batch) >= cap(e.wordQueue) {
			e.sendBatch(batch)
			batch = make([]string, 0, cap(e.wordQueue))
		}
	}

	// Channel is closed, send the last words
	e.sendBatch(batch)
	e.done <- true
}

// Finish closes the queue and
// waits for the last words to be send.
func (e *Endpoint) Finish() {
	close(e.wordQueue)
	<-e.done
}

// AddWord adds a word to the endpoint's queue.
//
// This function is thead-safe.
func (e *Endpoint) AddWord(word string) {
	e.wordQueue <- word
}

// New creates a new Endpoint worker with the specified destination URL and queue size.
func New(address string, queueSize int) (*Endpoint, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		conn:      conn,
		enc:       gob.NewEncoder(conn),
		wordQueue: make(chan string, queueSize),
		done:      make(chan bool),
	}, nil
}
