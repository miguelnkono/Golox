package errors

import "fmt"

var HadError bool

// perror for program error
func Perror(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where, message string) {
	err := fmt.Errorf("[Line %d] Error %s : %s\n", line, where, message)
	fmt.Print(err)
	HadError = true
}
