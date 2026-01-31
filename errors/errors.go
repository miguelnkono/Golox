package errors

import "fmt"

var HadError bool
var HadRuntimeError bool

// perror for program error
func Perror(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where, message string) {
	err := fmt.Errorf("[Line %d] Error %s : %s\n", line, where, message)
	fmt.Print(err)
	HadError = true
}

func ReportRuntimeError(err string) {
  fmt.Printf("Error: [%s]\n", err)
}
