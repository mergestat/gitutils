package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/mergestat/gitutils/clone"
)

func main() {
	args := os.Args[1:]
	err := clone.Exec(context.Background(), args[0], args[1])
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			log.Fatal(string(err.Stderr)+"\n", err)
		}
		log.Fatal(err)
	}
}
