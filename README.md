# gitutils

This is a Golang library for programmatically working with the `git` command (via the `os/exec` package).
When it comes to working with git repositories in go, there are some options:

  - [`go-git`](https://github.com/go-git/go-git) is a git implementation written in pure Go
  - [`git2go`](https://github.com/libgit2/git2go) are the Golang C-bindings to the `libgit2` project (requires CGO)
  - Shelling out to the `git` command and parsing results

This library uses the 3rd option (shelling out to `git`) and provides an abstraction layer to make using the output of various `git` subcommands easier.
