// Simple command line regular expression tester. Inspired by rubular.com.

package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

// COLOR is the ANSI color code for matched text. See https://en.wikipedia.org/wiki/ANSI_escape_code
const COLOR string = "34"

func usage() (output string) {
	output = "Simple CLI utility to test a regular expression against some input. Inspired by rubular.com.\n\n"
	output += "Full syntax: https://github.com/google/re2/wiki/Syntax\n\n"
	output += "Cheat sheet:\n"
	output += "[abc]\t\tA single character of: a, b, or c\t\t.\tAny single character\t\t\t\t\t(...)\tCapture everything enclosed\n"
	output += "[^abc]\t\tAny single character except: a, b, or c\t\t\\s\tAny whitespace character\t\t\t\t(a|b)\t a or b\n"
	output += "[a-z]\t\tAny single character in the range a-z\t\t\\S\tAny non-whitespace character\t\t\t\ta?\tZero or one of a\n"
	output += "[a-zA-Z]\tAny single character in the range a-z or A-Z\t\\d\tAny digit\t\t\t\t\t\ta*\tZero or more of a\n"
	output += "^\t\tStart of line\t\t\t\t\t\\D\tAny non-digit\t\t\t\t\t\ta+\tOne or more of a\n"
	output += "$\t\tEnd of line\t\t\t\t\t\\w\tAny word character (letter, number, underscore)\t\ta{3}\tExactly 3 of a\n"
	output += "\\A\t\tStart of string\t\t\t\t\t\\W\tAny non-word character\t\t\t\t\ta{3,}\t3 or more of a\n"
	output += "\\z\t\tEnd of string\t\t\t\t\t\\b\tAny word boundary\t\t\t\t\ta{3,6}\tBetween 3 and 6 of a\n\n"
	output += "Usage: " + path.Base(os.Args[0]) + " REGEX TESTCASE..."
	output += "\n"

	return
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println(usage())
		os.Exit(0)
	}

	re := regexp.MustCompile(os.Args[1])
	for _, testString := range os.Args[2:] {
		for _, match := range re.FindAll([]byte(testString), -1) {
			testString = strings.Replace(testString, string(match), "\033[1;"+COLOR+"m"+string(match)+"\033[0m", -1)
			fmt.Println(testString)
		}
	}
}
