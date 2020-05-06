package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/naoto0822/monkey-interpreter/pkg/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is monkey programing language!\n", user.Username)

	fmt.Printf("Feel free to type commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
