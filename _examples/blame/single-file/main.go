package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mergestat/gitutils/blame"
)

func main() {
	args := os.Args[1:]
	res, err := blame.Exec(context.Background(), args[0], args[1])
	if err != nil {
		log.Fatal(err)
	}

	var got strings.Builder
	for _, blame := range res {
		var linesInGroup string
		if blame.LinesInGroup != 0 {
			linesInGroup = fmt.Sprintf(" %d", blame.LinesInGroup)
		}
		got.WriteString(fmt.Sprintf("%s %d %d%s\n", blame.SHA, blame.OriginalLineNo, blame.FinalLineNo, linesInGroup))
		got.WriteString(fmt.Sprintf("author %s\n", blame.Author.Name))
		got.WriteString(fmt.Sprintf("author-mail <%s>\n", blame.Author.Email))
		got.WriteString(fmt.Sprintf("author-time %d\n", blame.Author.When.Unix()))
		got.WriteString(fmt.Sprintf("author-tz %s\n", blame.Author.When.Format("-0700")))
		got.WriteString(fmt.Sprintf("committer %s\n", blame.Committer.Name))
		got.WriteString(fmt.Sprintf("committer-mail <%s>\n", blame.Committer.Email))
		got.WriteString(fmt.Sprintf("committer-time %d\n", blame.Committer.When.Unix()))
		got.WriteString(fmt.Sprintf("committer-tz %s\n", blame.Committer.When.Format("-0700")))
		got.WriteString(fmt.Sprintf("summary %s\n", blame.Summary))
		if blame.Boundary {
			got.WriteString("boundary\n")
		}
		if blame.Previous != "" {
			got.WriteString(fmt.Sprintf("previous %s\n", blame.Previous))
		}
		got.WriteString(fmt.Sprintf("filename %s\n", blame.Filename))
		got.WriteString(fmt.Sprintf("\t%s\n", blame.Line))
	}

	fmt.Print(got.String())
}
