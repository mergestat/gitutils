package gitlog_test

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/mergestat/gitutils/gitlog"
)

var (
	repoPath string = "."
)

func init() {
	repoPath = os.Getenv("REPO_PATH")
	var err error
	repoPath, err = filepath.Abs(repoPath)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCount(t *testing.T) {
	iter, err := gitlog.Exec(context.Background(), repoPath)
	if err != nil {
		t.Fatal(err)
	}

	var count int
	for {
		commit, err := iter.Next()
		if err != nil {
			t.Fatal(err)
		}
		if commit == nil {
			break
		}
		count++
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.CommandContext(context.Background(), gitPath, "rev-list", "--count", "HEAD")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	wantCount, err := strconv.Atoi(strings.Trim(string(output), "\n"))
	if err != nil {
		t.Fatal(err)
	}

	if count != wantCount {
		t.Fatalf("mismatch in commit counts, got: %d want: %d", count, wantCount)
	}
}

func TestCountNoMerges(t *testing.T) {
	iter, err := gitlog.Exec(context.Background(), repoPath, gitlog.WithNoMerges(true))
	if err != nil {
		t.Fatal(err)
	}

	var count int
	for {
		commit, err := iter.Next()
		if err != nil {
			t.Fatal(err)
		}
		if commit == nil {
			break
		}
		count++
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.CommandContext(context.Background(), gitPath, "rev-list", "--count", "--no-merges", "HEAD")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	wantCount, err := strconv.Atoi(strings.Trim(string(output), "\n"))
	if err != nil {
		t.Fatal(err)
	}

	if count != wantCount {
		t.Fatalf("mismatch in commit counts, got: %d want: %d", count, wantCount)
	}
}

func TestLogOutput(t *testing.T) {
	iter, err := gitlog.Exec(context.Background(), repoPath)
	if err != nil {
		t.Fatal(err)
	}

	var got strings.Builder

	for {
		commit, err := iter.Next()
		if err != nil {
			t.Fatal(err)
		}
		if commit == nil {
			break
		}
		got.WriteString(commit.String() + "\n")
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.CommandContext(context.Background(), gitPath, "log", "--numstat", "--no-decorate", "-w", "--format=raw")
	cmd.Dir = repoPath

	w, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			t.Log(string(exitErr.Stderr))
		}
		t.Fatal(err)
	}

	want := string(w) + "\n"

	if string(want) != got.String() {
		t.Fatal("mismatch")
	}
}
