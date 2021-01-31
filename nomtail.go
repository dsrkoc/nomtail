package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	nextColor := NextColorIndexFn(Args.NoColor)

	fmt.Printf("getting job allocations from %s with job prefix '%s'\n", Args.Address, Args.JobPrefix)

	jobID, allocs, err := allocations()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if Args.Task == "" { // by default task id is the same as job id
		Args.Task = jobID
	}

	fmt.Println("Job Id:", jobID)
	fmt.Println("Number of allocations:", len(allocs))

	sigs := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(len(allocs))
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for _, allocID := range allocs {
		colIdx := nextColor()
		fmt.Println(Color(colIdx, "  allocation id:", allocID))

		go logs(colIdx, allocID, &wg)
	}

	go func() {
		sig := <-sigs
		fmt.Println("\nreceived signal:", sig)
		for i := 0; i < len(allocs); i++ {
			wg.Done() // artifically set WaitGroup counter to zero so app can exit
		}
	}()

	wg.Wait()
}
