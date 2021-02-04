package main

import (
	"flag"
	"fmt"
	"os"
)

type AppArgs struct {
	Address     string
	JobPrefix   string
	Task        string
	Type        string
	Namespace	string
	Follow      bool
	RunningOnly bool
	NoColor     bool
	Tail        int
}

var Args AppArgs

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s [OPTIONS] <job prefix>\n\nOptions:\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	nomadDefault := os.Getenv("NOMAD_ADDR")
	if nomadDefault == "" {
		nomadDefault = "http://localhost:4646"
	}

	flag.Usage = usage

	flag.StringVar(&Args.Address, "address", nomadDefault, "nomad's address")
	flag.StringVar(&Args.Task, "task", "", "Task id. Set if different from job id")
	flag.StringVar(&Args.Type, "type", "stdout", "stdout or stderr")
	flag.StringVar(&Args.Namespace, "namespace", "default", "specifies the target namespace")
	flag.BoolVar(&Args.Follow, "follow", false, "if set streams logs continually")
	flag.IntVar(&Args.Tail, "tail", 0, "shows the logs content with offsets relative to the end of the logs")
	flag.BoolVar(&Args.RunningOnly, "running-only", true, "if unset gets all allocations, not just the running ones")
	flag.BoolVar(&Args.NoColor, "no-color", false, "if set disables coloring of log lines")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	Args.JobPrefix = flag.Arg(0)
}
