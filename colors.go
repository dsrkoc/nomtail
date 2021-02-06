package main

import (
	"fmt"
	"math/rand"
)

const disableColors = -1

type decoration struct {
	Bold int
	Underline int
	Reversed int
}

// Decorations represent named values for various color decorations
var Decorations = decoration{Bold: 1, Underline: 4, Reversed: 7}

// NextColorIndexFn returns function that produces next color index
// to be used with the Color() function.
// The reson why an index producing function is returned, rather than
// index itself is that the returned function closes over the index
// state.
func NextColorIndexFn(noColor bool) func() int {
	if noColor {
		return func() int { return disableColors }
	}

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

// Color colorizes output using given color index
func Color(index int, args ...interface{}) string {
	if index == disableColors {
		return fmt.Sprint(args...)
	}
	template := fmt.Sprintf("\u001b[38;5;%dm%%s\u001b[0m", index)
	return fmt.Sprintf(template, fmt.Sprint(args...))
}

// Decor puts color decorations around given arguments
func Decor(index int, args ...interface{}) string {
	template := fmt.Sprintf("\u001b[%dm%%s\u001b[0m", index)
	return fmt.Sprintf(template, fmt.Sprint(args...))
}