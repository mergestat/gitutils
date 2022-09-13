package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mergestat/gitutils/lsfiles"
)

func main() {
	args := os.Args[1:]
	iter, err := lsfiles.Exec(context.Background(), args[0], lsfiles.WithFiles(args[1]))
	if err != nil {
		log.Fatal(err)
	}

	for {
		if file, err := iter.Next(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal(err)
		} else {
			fmt.Println(file)
		}
	}

}
