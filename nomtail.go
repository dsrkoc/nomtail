package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	nextColor := NextColorIndexFn(Args.NoColor)

	fmt.Printf("getting job allocations from %s with job prefix '%s'\n", Args.Address, Args.JobPrefix)

	jobID, allocs, err := allocations()
	if err != nil {
		fmt.Println(Decor(Decorations.Bold, "Error:"), err)
		os.Exit(1)
	}

	fmt.Println("Job Id:", jobID)
	fmt.Println("Number of allocations:", len(allocs))

	print := make(chan logEntry)
	stopPrint := make(chan bool)
	sigs := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(len(allocs))
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for _, alloc := range allocs {
		colIdx := nextColor()
		fmt.Println(" * allocation id:", Color(colIdx, alloc.ID), "("+alloc.State+")")

		go logs(colIdx, alloc.ID, print, &wg)
	}

	go printLog(500 * time.Millisecond, Args.Sort, print, stopPrint)

	go func() {
		sig := <-sigs
		fmt.Println("\nreceived signal:", sig)
		for i := 0; i < len(allocs); i++ {
			wg.Done() // artifically set WaitGroup counter to zero so app can exit
		}
	}()

	wg.Wait()
	stopPrint <- true
}
