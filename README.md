[![Go Reference](https://pkg.go.dev/badge/github.com/mergestat/gitutils.svg)](https://pkg.go.dev/github.com/mergestat/gitutils)
[![CI](https://github.com/mergestat/gitutils/actions/workflows/ci.yaml/badge.svg)](https://github.com/mergestat/gitutils/actions/workflows/ci.yaml)
[![Test Suite](https://github.com/mergestat/gitutils/actions/workflows/daily.yaml/badge.svg)](https://github.com/mergestat/gitutils/actions/workflows/daily.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mergestat/gitutils)](https://goreportcard.com/report/github.com/mergestat/gitutils)
[![codecov](https://codecov.io/gh/mergestat/gitutils/branch/main/graph/badge.svg?token=7KVW1U2LH7)](https://codecov.io/gh/mergestat/gitutils)

# gitutils

This is a Golang library for programmatically working with the `git` command (via the `os/exec` package).
In general, the options for working with git repositories in Go are:

  - [`go-git`](https://github.com/go-git/go-git) is a git implementation written in pure Go
  - [`git2go`](https://github.com/libgit2/git2go) are the Golang C-bindings to the `libgit2` project (requires CGO)
  - Shelling out to the `git` command (using the `os/exec` package) and parsing results

This library uses the 3rd option (shelling out to `git`) and provides an abstraction layer to simplify using the output of various `git` subcommands.

## Examples

### Cloning a repo

```golang
package main

import (
	"context"
	"log"
	"os"

	"github.com/mergestat/gitutils/clone"
)

func main() {
	err := clone.Exec(context.Background(), "https://github.com/mergestat/gitutils", "some-dir")
	if err != nil {
		log.Fatal(err)
	}
}
```

### Walking the Commit Log

```golang
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
	iter, err := gitlog.Exec(context.TODO(), "/path/to/some/local/repo", gitlog.WithStats(false))
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
```

See more examples in the [examples directory](https://github.com/mergestat/gitutils/tree/main/_examples).
