package main

import (
	"context"
	"log"
	"os"

	"github.com/mergestat/gitutils/clone"
)

func main() {
	args := os.Args[1:]
	configTest := map[string]string{"core.eol": "true"}
	err := clone.Exec(context.Background(), args[0], args[1], clone.WithConfig(configTest))
	if err != nil {
		log.Fatal(err)
	}
}
