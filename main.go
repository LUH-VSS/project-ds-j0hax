package main

import (
	"os"
	"sync"

	"github.com/LUH-VSS/project-ds-j0hax/AIO"
)

func main() {
	words := make(chan string)

	// For each file, start a ReaderWorker
	for _, file := range os.Args[1:] {
		go AIO.ReaderWorker(file, words)
	}

	// Start m WriterWorkers
	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go AIO.WriterWorker(&wg, "", words)
	}

	// Wait for Workers to finish...
	wg.Wait()

}
