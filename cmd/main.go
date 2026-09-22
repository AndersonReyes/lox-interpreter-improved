package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	lox "github.com/andersonreyes/lox/internal"
)

type builtInCommand = string

const (
	commandExit builtInCommand = ":exit"
)

func runFile(file string) {
	fmt.Printf("Running file %s\n", file)
	f, err := os.OpenFile(file, os.O_RDONLY, 0444)

	if err != nil {
		log.Fatalf("failed to open %s: %v", file, err)
	}

	reader := bufio.NewReader(f)
	scanner := lox.NewScanner(reader)

	err = scanner.Scan()
	if err != nil {
		log.Fatalf("error scanning file: %s\n", file)
	}
}

func runRepl() {
	fmt.Println("Welcome to the jlox repl!")
	inputReader := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")
	for inputReader.Scan() {
		if err := inputReader.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)

		} else {
			line := inputReader.Text()

			if line == commandExit {
				break
			}

			fmt.Println(line)
			// TODO: creating a new reader for each line seems inefficient? can we do better
			r := bufio.NewReader(strings.NewReader(line))
			scanner := lox.NewScanner(r)
			err := scanner.Scan()
			if err != nil {
				fmt.Printf("error interpreting line %s\n", line)
			}

			fmt.Print("> ")
		}
	}
	fmt.Println("goodbye!")
}

func main() {
	args := os.Args

	if len(args) < 1 {
		fmt.Println("Usage: jlox [script]")
		os.Exit(64) // command line usage error
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runRepl()
	}

}
