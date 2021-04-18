package main

import (
	"fmt"
	"sort"
	"time"
)

type logEntry struct {
	color int
	prefix string
	message string
}
type byLogEntry []logEntry

// Lexicographically sorting messages from the buffer
//

func (s byLogEntry) Len() int {
	return len(s)
}
func (s byLogEntry) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLogEntry) Less(i, j int) bool {
	return s[i].message < s[j].message
}

var buffer []logEntry = make([]logEntry, 0)

func printLog(collectDur time.Duration, sortBuffer bool, out <-chan logEntry, stop <-chan bool) {
	for {
		select {
		case entry := <-out:
			buffer = append(buffer, entry)
		case <-time.After(collectDur):
			if sortBuffer {
				sort.Sort(byLogEntry(buffer))
			}
			for _, elem := range(buffer) {
				shouldPrint := true
				for _, exclude := range Args.Excludes {
					shouldPrint = !exclude.MatchString(elem.message)
					if !shouldPrint { // one exclusion found is quite enough, thank you
						break
					}
				}

				if shouldPrint {
					fmt.Println(Color(elem.color, elem.prefix, elem.message))
				}
			}
			buffer = nil
		case <-stop:
			return
		}
	}
}