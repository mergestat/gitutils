package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mergestat/gitutils/gitlog"
)

func main() {
	args := os.Args[1:]

	iter, err := gitlog.Exec(context.Background(), args[0], gitlog.WithStats(false))
	if err != nil {
		log.Fatal(err)
	}

	for {
		if commit, err := iter.Next(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal(err)
		} else {
			fmt.Println(commit.SHA)
		}
	}
}
