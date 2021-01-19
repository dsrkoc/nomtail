package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Id struct { // both JSON structures contain "ID" field
	ID string
}

type alloc map[string]interface{}

func httpGet(url string) ([]byte, error) {
	resp, e1 := http.Get(url)
	if e1 != nil {
		return nil, e1
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func getJobs(nomadAddress string, jobPrefix string) ([]Id, error) {
	query := nomadAddress + "/v1/jobs?prefix=" + jobPrefix

	_jobs, e1 := httpGet(query)
	if e1 != nil {
		return nil, e1
	}

	var jobs []Id
	e2 := json.Unmarshal(_jobs, &jobs)

	return jobs, e2
}

func getAllocs(nomadAddress string, jobID string) ([]alloc, error) {
	query := nomadAddress + "/v1/job/%s/allocations"

	_allocs, e1 := httpGet(fmt.Sprintf(query, jobID))
	if e1 != nil {
		return nil, e1
	}

	var allocs []alloc
	e2 := json.Unmarshal(_allocs, &allocs)

	return allocs, e2
}

// func readState(alc alloc, jobID string) string {
// 	return alc["TaskStates"].(alloc)[jobID].(alloc)["State"].(string)
// }

// allocationIds returns a job indentifier and an array of that job's allocation identifiers.
func allocationIds(nomadAddress string, jobPrefix string) (string, []string, error) {
// It expects an address (e.g. address=http://localhost:4646) and job prefix

	// getting job identifier

	jobs, e1 := getJobs(nomadAddress, jobPrefix)
	if e1 != nil {
		return "", nil, e1
	}

	if len(jobs) > 1 {
		jobIds := make([]string, len(jobs))
		for i, job := range jobs {
			jobIds[i] = job.ID
		}
		joined := strings.Join(jobIds, ", ")
		return "", nil, fmt.Errorf("%d jobs found for given job prefix '%s' (%s)", len(jobs), jobPrefix, joined)
	}

	jobID := jobs[0].ID

	// getting list of allocation identifiers

	readState := func(alc alloc) string {
		return alc["TaskStates"].(alloc)[jobID].(alloc)["State"].(string)
	}

	allocs, e2 := getAllocs(nomadAddress, jobID)
	if e2 != nil {
		return "", nil, e2
	}

	allocIds := make([]string, len(allocs))
	for i, alloc := range allocs {
		allocIds[i] = alloc.ID
	}

	return jobID, allocIds, nil
}

func logs(color int, allocID string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(100 * time.Millisecond) // wait to allow main to print all the info before http request is sent

	url := fmt.Sprintf(
		"%s/v1/client/fs/logs/%s?follow=%t&type=%s&task=%s&origin=end&plain=true",
		Args.Nomad, allocID, Args.Follow, Args.Type, Args.Task)
	prefix := fmt.Sprintf("[%s] ", strings.Split(allocID, "-")[0]) // use only the first UUID segment

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting log for allocation "+Color(color, allocID)+":", err)
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			fmt.Println(Color(color, allocID), "done")
			return
		}
		if err != nil {
			if err == io.EOF {
				fmt.Println(Color(color, allocID), "done")
			} else {
				fmt.Println("Error reading log body for allocation "+Color(color, allocID)+":", err)
			}
			return
		}

		fmt.Println(Color(color, prefix, strings.TrimRight(string(line), "\n")))
	}
}
