package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mergestat/gitutils/lstree"
)

func main() {
	args := os.Args[1:]
	iter, err := lstree.Exec(context.Background(), args[0], args[1])
	if err != nil {
		log.Fatal(err)
	}

	for {
		if o, err := iter.Next(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal(err)
		} else {
			fmt.Println(o.String())
		}
	}
}
