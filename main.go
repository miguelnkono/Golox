package main

import (
	"bufio"
	"fmt"
	"golox/errors"
	"golox/parser"
	"golox/tools"
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
	// read the file content and executes it
	extension := filepath.Ext(path)
	if extension != ".golox" {
		log.Fatal("Script file must ends with the [.golox] extension")
	}
	source, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error: [%v]", err)
	}

	run(string(source))
	if errors.HadError {
		os.Exit(65)
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
	// execute the source code of the user
	// scan
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	// parse
	parser := parser.NewParser(tokens)
	expression, err := parser.Parse()

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	if expression != nil {
		astPrint := tools.AstPrinter{}
		fmt.Println(astPrint.Print(expression))
	}

	return nil
}
