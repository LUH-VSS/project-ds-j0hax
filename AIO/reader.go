package AIO

import (
	"bufio"
	"log"
	"os"
)

func ReaderWorker(file string, emitter chan<- string) {
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
	log.Printf("Finished reading %s\n", file)
}
