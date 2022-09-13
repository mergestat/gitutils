package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/mergestat/gitutils/blame"
	"github.com/mergestat/gitutils/lstree"
)

func main() {
	args := os.Args[1:]

	iter, err := lstree.Exec(context.Background(), args[0], "HEAD", lstree.WithRecurse(true))
	if err != nil {
		log.Fatal(err)
	}

	var objects []*lstree.Object
	for {
		if o, err := iter.Next(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			objects = append(objects, o)
		}
	}

	for _, o := range objects {
		if o.Type != "blob" {
			continue
		}
		res, err := blame.Exec(context.Background(), o.Path, &blame.Options{
			Directory: args[0],
		})
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Println("Error", err, o.Path, exitErr.Stderr)
			} else {
				fmt.Println("Error", err)
			}
			continue
		}

		fmt.Println("Blamed", len(res))
	}
}
