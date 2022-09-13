package lsfiles

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
)

type execOptions struct {
	Files            string
	NoEmptyDirectory bool
}

type Option func(o *execOptions)

// WithFiles corresponds to the first (and only) non-flag argument to git ls-files.
// It's a pattern for filtering the files to list. See <file> here: https://git-scm.com/docs/git-ls-files
func WithFiles(files string) Option {
	return func(o *execOptions) {
		o.Files = files
	}
}

// WithNoEmptyDirectory corresponds to the `--no-empty-directory` flag
// See here: https://git-scm.com/docs/git-ls-files#Documentation/git-ls-files.txt---no-empty-directory
func WithNoEmptyDirectory(NoEmptyDirectory bool) Option {
	return func(o *execOptions) {
		o.NoEmptyDirectory = NoEmptyDirectory
	}
}

type iterator struct {
	scanner *bufio.Scanner
}

// Next moves the iterator and returns the next file (or error).
// Iteration is complete when the error returned is io.EOF
func (i *iterator) Next() (string, error) {
	if next := i.scanner.Scan(); !next {
		if err := i.scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	} else {
		return i.scanner.Text(), nil
	}
}

// Exec runs `git ls-files` using the os/exec standard library package.
// It returns an iterator which can be used to retrieve a listing of files in a git repo.
// See here: https://git-scm.com/docs/git-ls-files
func Exec(ctx context.Context, repoPath string, options ...Option) (*iterator, error) {
	o := &execOptions{}
	for _, option := range options {
		option(o)
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("could not find git: %w", err)
	}

	args := []string{"ls-files"}

	if o.NoEmptyDirectory {
		args = append(args, "--no-empty-directory")
	}

	// NOTE: this has to be the last argument in the list
	if o.Files != "" {
		args = append(args, o.Files)
	}

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
