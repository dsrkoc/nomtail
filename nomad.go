package main

import (
	"bufio"
	"encoding/base64"
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

type alloc = map[string]interface{}

func httpGet(url string) ([]byte, error) {
	resp, e1 := http.Get(url)
	if e1 != nil {
		return nil, e1
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func getJobs(nomadAddress string, jobPrefix string, namespace string) ([]Id, error) {
	query := fmt.Sprintf("%s/v1/jobs?prefix=%s&namespace=%s", nomadAddress, jobPrefix, namespace)

	_jobs, e1 := httpGet(query)
	if e1 != nil {
		return nil, e1
	}

	var jobs []Id
	e2 := json.Unmarshal(_jobs, &jobs)

	return jobs, e2
}

func getAllocs(nomadAddress string, jobID string) ([]alloc, error) {
	query := fmt.Sprintf("%s/v1/job/%s/allocations", nomadAddress, jobID)

	_allocs, e1 := httpGet(query)
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

// allocations returns a job indentifier and an array of that job's allocation identifiers.
// It expects an address (e.g. address=http://localhost:4646) and job prefix
func allocations() (string, []string, error) {

	// getting job identifier

	jobs, e1 := getJobs(Args.Address, Args.JobPrefix, Args.Namespace)
	if e1 != nil {
		return "", nil, e1
	}

	if len(jobs) > 1 {
		jobIds := make([]string, len(jobs))
		for i, job := range jobs {
			jobIds[i] = job.ID
		}
		joined := strings.Join(jobIds, ", ")
		return "", nil, fmt.Errorf("%d jobs found for given job prefix '%s' (%s)", len(jobs), Args.JobPrefix, joined)
	}

	jobID := jobs[0].ID

	// getting list of allocation identifiers

	readState := func(alc alloc) string {
		return alc["TaskStates"].(alloc)[jobID].(alloc)["State"].(string)
	}

	allocs, e2 := getAllocs(Args.Address, jobID)
	if e2 != nil {
		return "", nil, e2
	}

	var allocIds []string
	for _, alloc := range allocs {
		id := alloc["ID"].(string)
		if Args.RunningOnly {
			if readState(alloc) == "running" {
				allocIds = append(allocIds, id)
			}
		} else {
			allocIds = append(allocIds, id)
		}
	}

	return jobID, allocIds, nil
}

func getLastLog(url string) ([]string, int, error) {
	type logJSON struct {
		Offset int
		Data   string
	}

	resp, err := httpGet(url)
	if err != nil {
		return nil, 0, err
	}

	logEntries := strings.Split(string(resp), "}{")

	// TODO perhaps it would be better if all logs could be joined instead of only last one taken
	var lastEntry string
	if len(logEntries) == 1 {
		lastEntry = logEntries[0]
	} else {
		lastEntry = "{" + logEntries[len(logEntries)-1]
	}

	var logEntry logJSON
	if err := json.Unmarshal([]byte(lastEntry), &logEntry); err != nil {
		return nil, 0, err
	}

	decodedLog, err := base64.StdEncoding.DecodeString(logEntry.Data)
	if err != nil {
		return nil, 0, err
	}
	return strings.Split(string(decodedLog), "\n"), logEntry.Offset, nil
}

func logs(color int, allocID string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(20 * time.Millisecond) // wait to allow main to print all the info before http request is sent

	prefix := fmt.Sprintf("[%s] ", strings.Split(allocID, "-")[0]) // use only the first UUID segment
	urlFirst := fmt.Sprintf("%s/v1/client/fs/logs/%s?type=%s&task=%s", Args.Address, allocID, Args.Type, Args.Task)

	url := urlFirst
	if Args.Tail > 0 {
		// get the first log batch so we can tail it
		lines, offset, err := getLastLog(urlFirst)
		if err != nil {
			fmt.Println("Error getting log for allocation "+Color(color, allocID)+":", err)
			return
		}
		url = fmt.Sprintf("%s&offset=%d", url, offset)

		for _, line := range lines[len(lines)-Args.Tail-1 : len(lines)-1] {
			fmt.Println((Color(color, prefix, line)))
		}

		if !Args.Follow { // it would seem that we're done
			fmt.Println(Color(color, allocID), "done")
			return
		}
	}

	urlRest := fmt.Sprintf("%s&follow=%t&plain=true", url, Args.Follow)

	resp, err := http.Get(urlRest)
	if err != nil {
		fmt.Println("Error getting log for allocation "+Color(color, allocID)+":", err)
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
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
