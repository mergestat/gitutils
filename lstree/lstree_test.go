package lstree_test

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mergestat/gitutils/lstree"
)

var (
	repoPath string = "."
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
}

func TestBasicOK(t *testing.T) {
	iter, err := lstree.Exec(context.Background(), repoPath, "HEAD", lstree.WithRecurse(true))
	if err != nil {
		t.Fatal(err)
	}

	var got strings.Builder
	for {
		o, err := iter.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatal(err)
		}
		got.WriteString(o.String() + "\n")
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.CommandContext(context.Background(), gitPath, "ls-tree", "HEAD", "-r")
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
		t.Fatal("mismatch")
	}
}
