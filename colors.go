package main

import (
	"fmt"
	"math/rand"
)

func NextIndexFn() func() int {
	i := 0

	return func() int {
		var next int

		if i > 13 {
			next = rand.Intn(210) + 20 // there's a chance of repeating numbers
		} else {
			i++
			next = i // numbers from 1 to 14
		}

		return next
	}
}

func Color(index int, args ...interface{}) string {
	template := fmt.Sprintf("\u001b[38;5;%dm%%s\u001b[0m", index)
	return fmt.Sprintf(template, fmt.Sprint(args...))
}
