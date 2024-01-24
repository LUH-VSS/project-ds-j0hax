package endpoint

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Endpoint struct {
	endpoint  *url.URL
	wordQueue chan string
}

// postWords sends a slice of words to the endpoint in JSON format.
//
// This function blocks until the bundle is successfully sent.
// In case of error, linear backoff is used.
func (e *Endpoint) postWords(words []string) {
	payload, err := json.Marshal(words)

	if err != nil {
		log.Print(err)
	}

	var body io.Reader = bytes.NewBuffer(payload)

	backoff := time.Second

	// Retry sending with a linear backoff until there is no error.
	// Once the bundle has been sent, reset it to a fresh slice.
	for {
		_, err = http.Post(e.endpoint.String(), "application/json", body)
		if err != nil {
			log.Print(err)
			log.Printf("Retrying in %s\n", backoff)
			time.Sleep(backoff)
			backoff += time.Second
		} else {
			break
		}
	}
}

// CollectWords collects words sent over the readers internal channel,
// and POSTs these in batches to the endpoint.
//
// A sync.WaitGroup is used to ensure the goroutine can completely finish.
func (e *Endpoint) CollectWords(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	batch := make([]string, 0, len(e.wordQueue))
	for word := range e.wordQueue {
		batch = append(batch, word)

		// Send all words in the batch once the slice is full
		if len(batch) == cap(batch) {
			e.postWords(batch)
			batch = make([]string, 0, len(e.wordQueue))
		}
	}

	// The channel has been closed, it is time to finish up.
	e.postWords(batch)
}

func (e *Endpoint) AddWord(word string) {
	e.wordQueue <- word
}

// Finish closes the internal queue channel.
//
// Assuming CollectWords has been called, this causes the remaining batch to be flushed.
func (e *Endpoint) Finish() {
	close(e.wordQueue)
}

// New creates a new Endpoint worker with the specified destination URL and queue size.
func New(destination string, queueSize uint64, wg *sync.WaitGroup) *Endpoint {
	url, err := url.Parse(destination)
	if err != nil {
		log.Panic(err)
	}

	return &Endpoint{
		endpoint:  url,
		wordQueue: make(chan string, queueSize),
	}
}
