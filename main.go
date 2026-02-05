package main

import (
	"bufio"
	"fmt"
	"golox/errors"
	"golox/parser"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) > 2 {
		log.Fatal("Usage: golox [script]\n")
	} else if len(os.Args) == 2 {
		// run script from file.
		runFile(os.Args[1])
	} else {
		// will lauch the prompt
		runPrompt()
	}
}

func runFile(path string) {
	extension := filepath.Ext(path)
	if extension != ".golox" {
		log.Fatal("Script file must ends with the [.golox] extension")
	}
	source, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error: [%v]", err)
	}

	if err := run(string(source)); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	if errors.HadError {
		os.Exit(65)
	}
	if errors.HadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	// read from the standard input
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		if err := run(line); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		errors.HadError = false
	}
}

func run(source string) error {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := parser.NewParser(tokens)
	expression, err := parser.Parse()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
  if expression == nil {
    return nil
  }

	intepreter := Interpreter{}
	intepreter.Interpret(expression)

	return nil
}
