package main

import (
	"bufio"
	"fmt"
	lox "github.com/andersonreyes/lox/internal"
	"os"
	"strings"
)

type builtInCommand = string

const (
	commandExit builtInCommand = ":exit"
)

func runFile(file string) {
	fmt.Printf("Running file %s\n", file)
}

func runRepl() {
	fmt.Println("Welcome to the jlox repl!\n")
	inputReader := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")
	for inputReader.Scan() {
		if err := inputReader.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)

		} else {
			line := inputReader.Text()

			if strings.Compare(line, commandExit) == 0 {
				break
			}

			fmt.Println(line)
			err := lox.ScanTokens(bufio.NewReader(strings.NewReader(line)))
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			fmt.Print("> ")
		}
	}
	fmt.Println("goodbye!")
}

func main() {
	args := os.Args

	if len(args) > 1 {
		fmt.Println("Usage: jlox [script]")
		os.Exit(64) // command line usage error
	} else if len(args) == 2 {
		runFile(args[2])
	} else {
		runRepl()
	}

}
