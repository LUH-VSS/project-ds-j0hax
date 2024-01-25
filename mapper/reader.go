package mapper

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/LUH-VSS/project-ds-j0hax/mapper/endpoint"
)

// Anything that is not a letter as defined by Unicode
var notLetters = regexp.MustCompile(`[^\p{L}]*`)

// Reader reads the specified files in parallel and uploads the JSON-encoded words to the endpoint
type Reader struct {
	files     []string
	endpoints []endpoint.Endpoint
	epwg      *sync.WaitGroup
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

	for _, ep := range r.endpoints {
		ep.Finish()
	}
	r.epwg.Wait()
	log.Printf("Finished POSTing to all mappers.")
}

// NewReader creates an instance of ReaderWorker
func NewReader(destinations []string, files []string) *Reader {

	var wg sync.WaitGroup

	endpoints := make([]endpoint.Endpoint, 0, len(destinations))
	for _, url := range destinations {
		ep := endpoint.New(url, 4096, &wg)
		go ep.CollectWords(&wg)
		endpoints = append(endpoints, *ep)
	}

	return &Reader{
		files:     files,
		endpoints: endpoints,
		epwg:      &wg,
	}
}
