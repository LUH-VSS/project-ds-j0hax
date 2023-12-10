package mapper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"
)

// Anything that is not a letter as defined by Unicode
var notLetters = regexp.MustCompile(`[^\p{L}]*`)

// Reader reads the specified files in parallel and uploads the JSON-encoded words to the endpoint
type Reader struct {
	files    []string
	endpoint *url.URL
}

// ProcessFile reads the file at the path in string and uploads the list to the endpoint
func (r *Reader) ProcessFile(file string) {
	log.Printf("Start reading %s\n", file)
	f, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	var wordList []string

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		rawWord := scanner.Text()
		word := notLetters.ReplaceAllString(rawWord, "") // remove punctuation
		if len(word) > 0 {
			wordList = append(wordList, word)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}
	log.Printf("Finished reading %s\n", file)
	r.postWords(wordList)
}

// postWords sends a slice of words to the endpoint in JSON format.
//
// This function blocks until the bundle is successfully sent.
// In case of error, linear backoff is used.
func (r *Reader) postWords(words []string) {
	log.Printf("Sending array of %d words\n", len(words))

	payload, err := json.Marshal(words)

	if err != nil {
		log.Print(err)
	}

	var body io.Reader = bytes.NewBuffer(payload)

	backoff := time.Second

	// Retry sending with a linear backoff until there is no error.
	// Once the bundle has been sent, reset it to a fresh slice.
	for {
		_, err = http.Post(r.endpoint.String(), "application/json", body)
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

// Run starts the ReaderWorker.
//
// It begins by reading files in parallel and
// POSTing these to its configured destination.
func (r *Reader) Run() {
	log.Printf("Mapping %d files\n", len(r.files))
	var wg sync.WaitGroup

	// Read each file in parallel
	for _, f := range r.files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			r.ProcessFile(f)
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
