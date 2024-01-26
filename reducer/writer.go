package reducer

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/LUH-VSS/project-ds-j0hax/lib"
)

type Writer struct {
	bindAddr        string
	outputFile      string
	incomingWords   chan string
	wordCounts      map[string]int64
	connectionGroup *sync.WaitGroup
}

// For each word recieved, check if it exists in the map and/or increment the value
func (w *Writer) countWords() {
	for word := range w.incomingWords {
		w.wordCounts[word] += 1
	}
}

// saveFile writes a tempfile containing the sorted words
// and their respective counts.
func (w *Writer) saveFile() {
	log.Printf("Sorting %d words\n", len(w.wordCounts))
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

// handleConnection handles a TCP connection.
//
// Received data is decoded and words are counted until EOF.
func (w *Writer) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer w.connectionGroup.Done()

	var data lib.Message
	dec := gob.NewDecoder(conn)

	for {
		err := dec.Decode(&data)
		if err != nil {
			if err == io.EOF {
				log.Printf("EOF from %s\n", conn.RemoteAddr())
				return
			}
			log.Panic(err)
		}
		for _, word := range data.Words {
			w.incomingWords <- word
		}
	}
}

// Run starts the reducer.
//
// It will accept incoming TCP connections, decode data and count words.
func (w *Writer) Run() {
	// Listen for a SIGINT and save file in that case.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Waiting for queues to finish...")
		w.connectionGroup.Wait()
		w.saveFile()
		os.Exit(0)
	}()

	go w.countWords()

	ln, err := net.Listen("tcp", w.bindAddr)
	if err != nil {
		panic(err)
	}

	log.Printf("Listening on %s\n", ln.Addr())
	log.Println("Press ^C (SIGINT) to save output file when done.")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		w.connectionGroup.Add(1)
		go w.handleConnection(conn)
	}
}

// NewWriter creates a reducer instance.
func NewWriter(addr, outputFile string) *Writer {
	gob.Register(lib.Message{})
	return &Writer{
		bindAddr:        addr,
		outputFile:      outputFile,
		incomingWords:   make(chan string),
		wordCounts:      make(map[string]int64),
		connectionGroup: &sync.WaitGroup{},
	}
}
