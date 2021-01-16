package main

/*
    /v1/jobs?prefix=<prefix> -> ID (job)
		/v1/job/<ID job>/allocations -> ID (allocation)
		/client/fs/logs/<ID (allocation)>
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Id struct { // both JSON structures contain "ID" field
	ID string
}

// getIds makes HTTP query to Nomad and returns array of Id types
func getIds(url string) ([]Id, error) {
	resp, e1 := http.Get(url)
	if e1 != nil {
		return nil, e1
	}
	defer resp.Body.Close()

	body, e2 := ioutil.ReadAll(resp.Body)
	if e2 != nil {
		return nil, e2
	}

	var ids []Id
	e3 := json.Unmarshal(body, &ids)
	return ids, e3
}

// allocationIds returns an array of allocation identifiers.
// It expects an address (e.g. address=http://locaohost:4646) and job prefix
func allocationIds(nomadAddress string, jobPrefix string) (string, []string, error) {
	queryJobs := nomadAddress + "/v1/jobs?prefix=" + jobPrefix
	queryAllocs := nomadAddress + "/v1/job/%s/allocations"

	// getting job identifier

	jobs, e1 := getIds(queryJobs)
	if e1 != nil {
		return "", nil, e1
	}

	if len(jobs) > 1 {
		jobIds := make([]string, len(jobs))
		for i, job := range jobs {
			jobIds[i] = job.ID
		}
		joined := strings.Join(jobIds, ", ")
		return "", nil, errors.New(fmt.Sprintf("%d jobs found for given job prefix '%s' (%s)", len(jobs), jobPrefix, joined))
	}

	jobId := jobs[0].ID

	// getting list of allocation identifiers

	allocs, e2 := getIds(fmt.Sprintf(queryAllocs, jobId))
	if e2 != nil {
		return "", nil, e2
	}

	allocIds := make([]string, len(allocs))
	for i, alloc := range allocs {
		allocIds[i] = alloc.ID
	}

	return jobId, allocIds, nil
}

// main -----------------------

func main() {
	nextColor := NextIndexFn()
	args := processCmdLineArgs()

	fmt.Printf("- getting job allocations from %s with job prefix '%s'\n", args.Nomad, args.JobPrefix)

	jobId, allocs, err := allocationIds(args.Nomad, args.JobPrefix)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if args.Task == "" { // by default task id is the same as job id
		args.Task = jobId
	}

	fmt.Println("Job Id:", jobId)
	fmt.Println("Number of allocations:", len(allocs))
	for _, allocId := range allocs {
		colIdx := nextColor()
		fmt.Println(Color(colIdx, "  allocation id:", allocId))
	}

	fmt.Println("\n<== Done")

}
