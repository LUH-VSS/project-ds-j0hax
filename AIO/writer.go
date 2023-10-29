package AIO

import (
	"log"
	"sync"
)

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
