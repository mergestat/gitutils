package main

import (
	"context"
	"log"
	"os"

	"github.com/mergestat/gitutils/clone"
)

func main() {
	args := os.Args[1:]
	serverOpt := []string{"asdasd", "adasdad"}
	err := clone.Exec(context.Background(), args[0], args[1], clone.WithServerOption(serverOpt))
	if err != nil {
		log.Fatal(err)
	}
}
