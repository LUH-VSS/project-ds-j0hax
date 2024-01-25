package mapper

import (
	"bufio"
	"encoding/gob"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/LUH-VSS/project-ds-j0hax/lib"
	"github.com/LUH-VSS/project-ds-j0hax/mapper/endpoint"
)

// Anything that is not a letter as defined by Unicode
var notLetters = regexp.MustCompile(`[^\p{L}]*`)

// Reader reads the specified files in parallel and uploads the JSON-encoded words to the endpoint
type Reader struct {
	files     []string
	endpoints []endpoint.Endpoint
}

// ProcessFile reads the file at the path in string and uploads the list to the endpoint
func (r *Reader) ProcessFile(file string) {
	log.Printf("Start reading %s\n", file)
	f, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		rawWord := scanner.Text()
		word := notLetters.ReplaceAllString(rawWord, "") // remove punctuation

		r.Map(word)
	}
	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}
	log.Printf("Finished reading %s\n", file)
}

// Hash hashes a word according to Dan Bernstein's DJB2 algorithm.
func Hash(word string) uint64 {
	hash := uint64(5381)
	for _, r := range word {
		hash = ((hash << 5) + hash) ^ uint64(r)
	}
	return hash
}

// Map deterministically chooses an
// endpoint for a word, and adds it to its queue.
func (r *Reader) Map(word string) {
	// Pick an endpoint for this word
	index := Hash(word) % uint64(len(r.endpoints))
	r.endpoints[index].AddWord(word)
}

// Run starts the ReaderWorker.
//
// It begins by reading files in parallel and
// POSTing these to its configured destination.
func (r *Reader) Run() {
	for _, e := range r.endpoints {
		go e.Run()
	}

	log.Printf("Mapping %d files\n", len(r.files))
	var fileGroup sync.WaitGroup
	// Read each file in parallel
	for _, f := range r.files {
		fileGroup.Add(1)
		go func(f string) {
			defer fileGroup.Done()
			r.ProcessFile(f)
		}(f)
	}

	fileGroup.Wait()
	log.Printf("Finished reading all files.")

	// Flush remaining batched words
	for _, e := range r.endpoints {
		e.Finish()
	}
	log.Println("Finished sending to all endpoints.")
}

// NewReader creates an instance of ReaderWorker
func NewReader(destinations []string, files []string) *Reader {
	gob.Register(lib.Message{})
	endpoints := make([]endpoint.Endpoint, 0, len(destinations))
	for _, url := range destinations {
		ep, err := endpoint.New(url, 4096)
		if err != nil {
			log.Printf("Skipping Reducer: %e\n", err)
		}

		endpoints = append(endpoints, *ep)
	}

	return &Reader{
		files:     files,
		endpoints: endpoints,
	}
}
