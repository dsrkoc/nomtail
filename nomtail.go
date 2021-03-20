package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(0) // log is used as an easier way to print messages other than services' logs to stderr

	nextColor := NextColorIndexFn(Args.NoColor)

	log.Printf("getting job allocations from %s with job prefix '%s'\n", Args.Address, Args.JobPrefix)

	jobID, allocs, err := allocations()
	if err != nil {
		log.Fatalln(Decor(Decorations.Bold, "Error:"), err)
	}

	log.Println("Job Id:", jobID)
	log.Println("Number of allocations:", len(allocs))

	// messages buffer gets emptied to stdout periodically,
	// every collectMsgsDur milliseconds
	collectMsgsDur := 500 * time.Millisecond

	print := make(chan logEntry)
	stopPrint := make(chan bool)
	sigs := make(chan os.Signal, 1)

	var wg sync.WaitGroup
	var loggersCount int32

	loggerDoneFn := func () {
		wg.Done()
		atomic.AddInt32(&loggersCount, -1)
	}

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for _, alloc := range allocs {
		colIdx := nextColor()
		log.Println(" * allocation id:", Color(colIdx, alloc.ID), "("+alloc.State+")")

		wg.Add(1)
		loggersCount++ // safe to do this without atomic since loggers sleep a while

		go logs(colIdx, alloc.ID, collectMsgsDur, print, loggerDoneFn)
	}

	go printLog(collectMsgsDur, Args.Sort, print, stopPrint)

	go func() {
		sig := <-sigs
		log.Println("\nreceived signal:", sig)

		loggersAlive := atomic.LoadInt32(&loggersCount)
		for i := 0; i < int(loggersAlive); i++ {
			wg.Done() // artifically set WaitGroup counter to zero so app can exit
		}
	}()

	wg.Wait()
	stopPrint <- true
}
