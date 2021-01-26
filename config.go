package main

import (
	"flag"
	"os"
)

type AppArgs struct {
	Address     string
	JobPrefix   string
	Task        string
	Type        string
	Follow      bool
	RunningOnly bool
	NoColor     bool
	// Tail        int
}

var Args AppArgs

func init() {
	nomadDefault := os.Getenv("NOMAD_ADDR")
	if nomadDefault == "" {
		nomadDefault = "http://localhost:4646"
	}

	flag.StringVar(&Args.Address, "address", nomadDefault, "nomad's address")
	flag.StringVar(&Args.JobPrefix, "job-prefix", "unknown", "job prefix (should uniquely identify a job)")
	flag.StringVar(&Args.Task, "task", "", "Task id. Set if different from job id")
	flag.StringVar(&Args.Type, "type", "stdout", "stdout or stderr")
	flag.BoolVar(&Args.Follow, "follow", false, "if set streams logs continually")
	// flag.IntVar(&Args.Tail, "tail", 10, "shows the logs content with offsets relative to the end of the logs")
	flag.BoolVar(&Args.RunningOnly, "running-only", true, "if unset gets all allocations, not just the running ones")
	flag.BoolVar(&Args.NoColor, "no-color", false, "if set disables coloring of log lines")

	flag.Parse()
}
