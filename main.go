package main

import (
	"bufio"
	"fmt"
	"golox/errors"
	"golox/parser"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	switch len(os.Args) {
	case 1:
		runPrompt()
	case 2:
		runFile(os.Args[1])
	default:
		fmt.Fprintln(os.Stderr, "Usage: golox [script]")
		os.Exit(64)
	}
}

var interpreter = NewInterpreter()

func runFile(path string) {
	extension := filepath.Ext(path)
	if extension != ".golox" && extension != ".lox" {
		log.Fatal("Script file must end with .golox or .lox extension")
	}

	source, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	run(string(source))

	// Exit with appropriate error code
	if errors.HadError {
		os.Exit(65)
	}
	if errors.HadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("GoLox REPL - Type 'exit' or 'quit' to quit")

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

		// Allow user to exit gracefully
		if line == "exit" || line == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		if line == "" {
			continue
		}

		run(line)

		// Reset error flag in REPL mode so user can continue
		errors.HadError = false
		errors.HadRuntimeError = false
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

func run(source string) {
	// Scanning
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	// Stop if there were scan errors
	if errors.HadError {
		return
	}

	// Parsing
	parser := parser.NewParser(tokens)
	statements, err := parser.Parse()

	if err != nil {
		fmt.Println(err)
		errors.HadError = true
		return
	}

	// Interpreting
	interpreter.Interpret(statements)
}
