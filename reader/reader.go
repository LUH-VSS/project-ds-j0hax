package AIO

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

type ReaderWorker struct {
	files    []string
	endpoint *url.URL
}

func readFile(file string, emitter chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		emitter <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}
}

func (r *ReaderWorker) sendWords(emitter <-chan string) {
	for word := range emitter {
		dat := &data.Word{
			Word: word,
		}

		payload, err := json.Marshal(dat)
		if err != nil {
			log.Print(err)
		}

		var body io.Reader = bytes.NewBuffer(payload)
		http.Post(r.endpoint.String(), "application/json", body)
	}
}

// Run starts the ReaderWorker.
// It begins by reading files in parallel and
// POSTing these to its configured destination.
func (r *ReaderWorker) Run() {
	var wg sync.WaitGroup
	buffsize := len(r.files)
	words := make(chan string, buffsize)

	go r.sendWords(words)

	for _, f := range r.files {
		wg.Add(1)
		go readFile(f, words, &wg)
	}

	wg.Wait()
	log.Printf("Transmitted all files")
}

// NewReader creates an instance of ReaderWorker
func NewReader(destination string, files []string) *ReaderWorker {
	url, err := url.Parse(destination)
	if err != nil {
		log.Panic(err)
	}

	return &ReaderWorker{
		files:    files,
		endpoint: url,
	}
}
