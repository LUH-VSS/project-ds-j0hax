package reducer

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/LUH-VSS/project-ds-j0hax/data"
)

type Writer struct {
	bindAddr string
	pattern  string
}

func WriterWorker(wg *sync.WaitGroup, file string, collector <-chan string) {
	defer wg.Done()

	// Internal hashmap for storing word counts
	m := make(map[string]int)

	// For each word recieved, check if it exists in the map and/or increment the value
	for word := range collector {
		val, ok := m[word]
		if ok {
			m[word] = val + 1
		} else {
			m[word] = 1
		}

		log.Printf("%s = %d\n", word, m[word])
	}
}

func test(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t data.Word
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Printf("%#v\n", t)
}

func (w *Writer) Run() {
	log.Printf("Listening on %s%s\n", w.bindAddr, w.pattern)
	http.HandleFunc(w.pattern, test)
	log.Fatal(http.ListenAndServe(w.bindAddr, nil))
}

func NewWriter(addr, pattern string) *Writer {
	return &Writer{
		bindAddr: addr,
		pattern:  pattern,
	}
}
