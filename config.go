package main

import (
	"flag"
	"os"
)

type Args struct {
	Nomad     string
	JobPrefix string
	Task      string
	Type      string
	Follow    bool
	Tail      int
}

func processCmdLineArgs() Args {
	nomadDefault := os.Getenv("NOMAD_ADDR")
	if nomadDefault == "" {
		nomadDefault = "http://localhost:4646"
	}

	nomad := flag.String("nomad", nomadDefault, "nomad URI")
	jobPrefix := flag.String("job-prefix", "unknown", "job prefix (should uniquely identify a job)")
	task := flag.String("task", "", "Task id. Set if different from job id")
	typ := flag.String("type", "stdout", "stdout or stderr")
	follow := flag.Bool("no-follow", false, "if set pulls logs and stops")
	tail := flag.Int("tail", 10, "shows the logs content with offsets relative to the end of the logs")

	flag.Parse()

	return Args{
		Nomad:     *nomad,
		JobPrefix: *jobPrefix,
		Task:      *task,
		Type:      *typ,
		Follow:    *follow,
		Tail:      *tail,
	}
}
