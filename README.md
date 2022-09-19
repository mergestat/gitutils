# gitutils

This is a Golang library for programmatically working with the `git` command (via the `os/exec` package).
There are some options for working with git repositories in go:

  - [`go-git`](https://github.com/go-git/go-git) is a git implementation written in pure Go
  - [`git2go`](https://github.com/libgit2/git2go) are the Golang C-bindings to the `libgit2` project (requires CGO)
  - Shelling out to the `git` command (using the `os/exec` package) and parsing results

This library uses the 3rd option (shelling out to `git`) and provides an abstraction layer to make using the output of various `git` subcommands easier.

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

more examples on the way...
