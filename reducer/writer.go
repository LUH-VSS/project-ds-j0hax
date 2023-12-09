package reducer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"
)

type Writer struct {
	bindAddr      string
	pattern       string
	outputFile    string
	incomingWords chan string
	wordCounts    map[string]int64
}

// For each word recieved, check if it exists in the map and/or increment the value
func (w *Writer) countWords() {
	for word := range w.incomingWords {
		w.wordCounts[word] += 1
	}
}

func (w *Writer) saveFile() {
	keys := make([]string, 0, len(w.wordCounts))
	for k := range w.wordCounts {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	file, err := os.CreateTemp("", "excercise")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, k := range keys {
		fmt.Fprintf(file, "%s %d\n", k, w.wordCounts[k])
	}

	log.Printf("Saved to %s\n", file.Name())
}

func (w *Writer) handleWords(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	var wrd string
	err := decoder.Decode(&wrd)
	if err != nil {
		panic(err)
	}

	w.incomingWords <- wrd
}

func (w *Writer) Run() {
	// Listen for a SIGINT and save file in that case.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		w.saveFile()
		os.Exit(0)
	}()

	go w.countWords()

	log.Printf("Listening on %s under %s\n", w.bindAddr, w.pattern)
	log.Println("Press ^C (SIGINT) to save output file when done.")
	http.HandleFunc(w.pattern, w.handleWords)
	log.Fatal(http.ListenAndServe(w.bindAddr, nil))
}

func NewWriter(addr, pattern, outputFile string) *Writer {
	return &Writer{
		bindAddr:      addr,
		pattern:       pattern,
		outputFile:    outputFile,
		incomingWords: make(chan string),
		wordCounts:    make(map[string]int64),
	}
}
