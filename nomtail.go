package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// main -----------------------

func main() {
	nextColor := NextIndexFn()

	fmt.Printf("getting job allocations from %s with job prefix '%s'\n", Args.Nomad, Args.JobPrefix)

	jobId, allocs, err := allocationIds(Args.Nomad, Args.JobPrefix)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if Args.Task == "" { // by default task id is the same as job id
		Args.Task = jobId
	}

	fmt.Println("Job Id:", jobId)
	fmt.Println("Number of allocations:", len(allocs))

	sigs := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(len(allocs))
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for _, allocId := range allocs {
		colIdx := nextColor()
		fmt.Println(Color(colIdx, "  allocation id:", allocId))

		go logs(colIdx, allocId, &wg)
	}

	go func() {
		sig := <-sigs
		fmt.Println("\nreceived signal:", sig)
		for i := 0; i < len(allocs); i++ {
			wg.Done() // artifically set WaitGroup counter to zero so app can exit
		}
	}()

	wg.Wait()
	fmt.Println("<== Done")

}
