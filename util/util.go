package util

import (
	"fmt"
	"log"
	"strings"
)

// ShowTest ...
func ShowTest(text string) {
	fmt.Println(strings.Repeat("- -", 20))
	fmt.Println(`
	╦═╗┌─┐┌┐┌┌┬┐┌─┐┌┬┐  ┌┬┐┌─┐┌─┐┌┬┐┌─┐ 
	╠╦╝├─┤│││ │││ ││││   │ ├┤ └─┐ │ └─┐ 
	╩╚═┴ ┴┘└┘─┴┘└─┘┴ ┴   ┴ └─┘└─┘ ┴ └─┘o
	STDOUT: contents that may match your internal profile wall
	`)
	fmt.Println(strings.Repeat("- -", 20))
	fmt.Println(text)
	fmt.Println(strings.Repeat("- -", 20))

}


// Log .. 
func Log(format string, v ...interface{})  {
	log.Printf(format, v...)
}