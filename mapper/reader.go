package mapper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"io"

	"github.com/LUH-VSS/project-ds-j0hax/data"
)

// Reader reads the specified files in parallel and uploads the JSON-encoded words to the endpoint
type Reader struct {
	files    []string
	endpoint *url.URL
}

// ReadFile reads the file at the path in string and sends the words contained in it to the emitter channel. wg.Done() is called when the function exits.
func ReadFile(file string, emitter chan<- data.Word) {
	log.Printf("Start reading %s\n", file)
	f, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		emitter <- data.Word{
			OriginFile: file,
			Word:       scanner.Text(),
		}
	}
	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}
	log.Printf("Finished reading %s\n", file)
}

// sendWords
func (r *Reader) sendWords(emitter <-chan data.Word) {
	for word := range emitter {
		payload, err := json.Marshal(word)
		if err != nil {
			log.Print(err)
		}

		var body io.Reader = bytes.NewBuffer(payload)
		_, err = http.Post(r.endpoint.String(), "application/json", body)
		if err != nil {
			log.Print(err)
		}
	}
}

// Run starts the ReaderWorker.
//
// It begins by reading files in parallel and
// POSTing these to its configured destination.
func (r *Reader) Run() {
	var wg sync.WaitGroup
	buffsize := len(r.files)
	words := make(chan data.Word, buffsize)

	// Start the goroutine that listens on the channel and sends the words
	go r.sendWords(words)

	// Read each file in parallel
	for _, f := range r.files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			ReadFile(f, words)
		}(f)
	}
	wg.Wait()
}

// NewReader creates an instance of ReaderWorker
func NewReader(destination string, files []string) *Reader {
	url, err := url.Parse(destination)
	if err != nil {
		log.Panic(err)
	}

	return &Reader{
		files:    files,
		endpoint: url,
	}
}
