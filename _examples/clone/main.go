package main

import (
	"context"
	"log"
	"os"

	"github.com/mergestat/gitutils/clone"
)

func main() {
	args := os.Args[1:]
	err := clone.Exec(context.Background(), args[0], args[1], clone.WithBare(true))
	if err != nil {
		log.Fatal(err)
	}
}
