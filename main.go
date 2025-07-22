package main

import (
	"os"
	"os/user"
	"fmt"
	"log"	

	"monkey/repl"
)

func main() {
	// Gets the current OS session's user's name
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hello, %s! This is the Monkey Programming Language REPL!\n", user.Username)
	fmt.Printf("Enter commands after the monkey prompt.\n\n")

	// start REPL (language "shell")
	repl.Start(os.Stdin, os.Stdout)
}
