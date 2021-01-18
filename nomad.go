package main

import (
	"bufio"
	"encoding/json"
	"errors"
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

// allocationIds returns a job indentifier and an array of that job's allocation identifiers.
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

func logs(color int, allocId string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(100 * time.Millisecond) // wait to allow main to print all the info before http request is sent

	url := fmt.Sprintf(
		"%s/v1/client/fs/logs/%s?follow=%t&type=%s&task=%s&origin=end&plain=true",
		Args.Nomad, allocId, Args.Follow, Args.Type, Args.Task)
	prefix := fmt.Sprintf("[%s] ", strings.Split(allocId, "-")[0]) // use only the first UUID segment

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting log for allocation "+Color(color, allocId)+":", err)
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			fmt.Println(Color(color, allocId), "done")
			return
		}
		if err != nil {
			if err == io.EOF {
				fmt.Println(Color(color, allocId), "done")
			} else {
				fmt.Println("Error reading log body for allocation "+Color(color, allocId)+":", err)
			}
			return
		}

		fmt.Println(Color(color, prefix, strings.TrimRight(string(line), "\n")))
	}
}
