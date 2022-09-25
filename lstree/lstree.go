// Package lstree shells out to git ls-tree https://git-scm.com/docs/git-ls-tree
package lstree

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Mode string

const (
	NormalFile     Mode = "100644"
	ExecutableFile Mode = "100755"
	SymbolicLink   Mode = "120000"
	Tree           Mode = "040000"
	Submodule      Mode = "160000"
)

type execOptions struct {
	Recurse bool
}

type Option func(o *execOptions)

// WithRecurse corresponds to the -r flag
// https://www.git-scm.com/docs/git-ls-tree#Documentation/git-ls-tree.txt--r
func WithRecurse(recurse bool) Option {
	return func(o *execOptions) {
		o.Recurse = recurse
	}
}

type iterator struct {
	scanner *bufio.Scanner
}

type Object struct {
	Mode string
	Type string
	Hash string
	Path string
}

// modeFromString returns a Mode from a string representation of a git object mode.
func modeFromString(s string) Mode {
	switch s {
	case "100644":
		return NormalFile
	case "100755":
		return ExecutableFile
	case "120000":
		return SymbolicLink
	case "040000":
		return Tree
	case "160000":
		return Submodule
	default:
		return ""
	}
}

// objectFromOutputLine parses a single line in the default git ls-tree output format
// and returns an Object struct. See here: https://www.git-scm.com/docs/git-ls-tree#_output_format
func objectFromOutputLine(line string) *Object {
	s := strings.SplitN(line, " ", 3)
	s2 := strings.Split(s[2], "\t")
	return &Object{
		Mode: string(modeFromString(s[0])),
		Type: s[1],
		Hash: s2[0],
		Path: s2[1],
	}
}

// String returns an *Object in the same format as a single line of the default
// git ls-tree output. See here: https://www.git-scm.com/docs/git-ls-tree#_output_format
func (o *Object) String() string {
	return fmt.Sprintf("%s %s %s\t%s", o.Mode, o.Type, o.Hash, o.Path)
}

// Next moves the iterator and returns the next object (or error).
// Iteration is complete when the error returned is io.EOF
func (i *iterator) Next() (*Object, error) {
	if next := i.scanner.Scan(); !next {
		if err := i.scanner.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	} else {
		return objectFromOutputLine(i.scanner.Text()), nil
	}
}

// Exec runs `git ls-tree` using the os/exec standard library package.
// It returns an iterator which can be used to retrieve the contents of a tree-ish
// See here: https://www.git-scm.com/docs/git-ls-tree
func Exec(ctx context.Context, repoPath, treeish string, options ...Option) (*iterator, error) {
	o := &execOptions{}
	for _, option := range options {
		option(o)
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("could not find git: %w", err)
	}

	args := []string{"ls-tree"}

	if o.Recurse {
		args = append(args, "-r")
	}

	args = append(args, treeish)

	cmd := exec.CommandContext(ctx, gitPath, args...)
	cmd.Dir = repoPath

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	iter := &iterator{
		scanner: bufio.NewScanner(stdout),
	}

	return iter, nil
}
