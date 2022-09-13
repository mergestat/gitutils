package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mergestat/gitutils/blame"
)

func main() {
	args := os.Args[1:]
	fmt.Println(args)
	res, err := blame.Exec(context.Background(), args[1], &blame.Options{
		Directory: args[0],
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
