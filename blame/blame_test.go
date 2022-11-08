package blame_test

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mergestat/gitutils/blame"
)

var (
	repoPath string = "../" // blame can run from a subdir of a repo
	filePath string = "README.md"
)

func init() {
	if r := os.Getenv("REPO_PATH"); r != "" {
		repoPath = os.Getenv("REPO_PATH")
	}
	var err error
	repoPath, err = filepath.Abs(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	if p := os.Getenv("BLAME_FILE_PATH"); p != "" {
		filePath = p
	}
}

func TestBlameOutput(t *testing.T) {
	adjustedBufferSize := bufio.MaxScanTokenSize * 5
	res, err := blame.Exec(context.Background(), repoPath, filePath, blame.WithScannerBuffer(make([]byte, adjustedBufferSize), adjustedBufferSize))
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			t.Log(string(exitErr.Stderr))
		}
		t.Fatal(err)
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

	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.CommandContext(context.Background(), gitPath, "blame", "--line-porcelain", filePath)
	cmd.Dir = repoPath

	w, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			t.Log(string(exitErr.Stderr))
		}
		t.Fatal(err)
	}

	want := string(w)

	if string(want) != got.String() {
		fmt.Println(len(string(want)))
		fmt.Println(len(got.String()))
		t.Fatal("mismatch")
	}
}
