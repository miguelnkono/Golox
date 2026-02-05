package errors

import (
	"fmt"
	"io"
	"os"
)

var HadError bool
var HadRuntimeError bool

type ErrorReporter struct {
	writer io.Writer
}

func NewErrorReporter(w io.Writer) *ErrorReporter {
	if w == nil {
		w = os.Stderr
	}
	return &ErrorReporter{writer: w}
}

// Default reporter writes to stderr
var DefaultReporter = NewErrorReporter(os.Stderr)

// ScanError reports errors from the scanner
func (e *ErrorReporter) ScanError(line int, message string) {
	e.report(line, "", message, "ScanError")
	HadError = true
}

// ParseError reports errors from the parser
func (e *ErrorReporter) ParseError(line int, where, message string) {
	e.report(line, where, message, "ParseError")
	HadError = true
}

// RuntimeError reports errors from the interpreter
func (e *ErrorReporter) RuntimeError(line int, message string) {
	e.report(line, "", message, "RuntimeError")
	HadRuntimeError = true
}

func (e *ErrorReporter) report(line int, where, message, errorType string) {
	whereStr := ""
	if where != "" {
		whereStr = fmt.Sprintf(" at '%s'", where)
	}

	fmt.Fprintf(e.writer, "[line %d] %s%s: %s\n", line, errorType, whereStr, message)
}

// Legacy functions for backward compatibility
func Perror(line int, message string) {
	DefaultReporter.ScanError(line, message)
}

func Report(line int, where, message string) {
	DefaultReporter.ParseError(line, where, message)
}

func ReportRuntimeError(err string) {
	fmt.Fprintf(DefaultReporter.writer, "RuntimeError: %s\n", err)
	HadRuntimeError = true
}
