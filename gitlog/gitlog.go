// Package gitlog ...
package gitlog

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const isoDataFmtStr = "2006-01-02T15:04:05-07:00"

// line prefixes used in the format string
const (
	commitHashPrefix     = "_H:"
	treeHashPrefix       = "_T:"
	parentHashesPrefix   = "_P:"
	authorNamePrefix     = "_aN:"
	authorEmailPrefix    = "_aE:"
	authorDatePrefix     = "_aI:"
	committerNamePrefix  = "_cN:"
	committerEmailPrefix = "_cE:"
	committerDatePrefix  = "_cI:"
	commitBodyPrefix     = "_B:"
)

// buildFormatString constructs a format string to pass to `git log`
func buildFormatString() string {
	var b strings.Builder
	b.WriteString(commitHashPrefix + "%H%n")
	b.WriteString(treeHashPrefix + "%T%n")
	b.WriteString(parentHashesPrefix + "%P%n")

	b.WriteString(authorNamePrefix + "%aN%n")
	b.WriteString(authorEmailPrefix + "%aE%n")
	b.WriteString(authorDatePrefix + "%aI%n")

	b.WriteString(committerNamePrefix + "%cN%n")
	b.WriteString(committerEmailPrefix + "%cE%n")
	b.WriteString(committerDatePrefix + "%cI%n")

	b.WriteString(commitBodyPrefix + "%B%n%x00")

	return b.String()
}

// Commit represents a commit parsed from git log
type Commit struct {
	SHA                string
	Tree               string
	Parents            []string
	Author             Event
	Committer          Event
	Message            string
	Stats              []Stat
	hasTrailingNewline bool
}

// Event represents the who and when of a commit event
type Event struct {
	Name  string
	Email string
	When  time.Time
}

// Stat holds the diff stat of a file
type Stat struct {
	FilePath  string
	Additions int
	Deletions int
}

type execOptions struct {
	NoMerges     bool
	IncludeStats bool
	FileFilter   string
	Order        CommitOrder
	FirstParent  bool
	M            bool
	Stats        bool
}

type CommitOrder string

const (
	DateOrder       CommitOrder = "--date-order"
	AuthorDateOrder CommitOrder = "--author-data-order"
	TopoOrder       CommitOrder = "--topo-order"
	Reverse         CommitOrder = "--reverse"
)

type Option func(o *execOptions)

func WithNoMerges(noMerges bool) Option {
	return func(o *execOptions) {
		o.NoMerges = noMerges
	}
}

func WithIncludeStats(includeStats bool) Option {
	return func(o *execOptions) {
		o.IncludeStats = includeStats
	}
}

func WithFileFilter(fileFilter string) Option {
	return func(o *execOptions) {
		o.FileFilter = fileFilter
	}
}

func WithCommitOrder(order CommitOrder) Option {
	return func(o *execOptions) {
		o.Order = order
	}
}

func WithFirstParent(firstParent bool) Option {
	return func(o *execOptions) {
		o.FirstParent = firstParent
	}
}

func WithM(m bool) Option {
	return func(o *execOptions) {
		o.M = m
	}
}

func WithStats(stats bool) Option {
	return func(o *execOptions) {
		o.Stats = stats
	}
}

type commitIterator struct {
	// scanner is a Scanner produced from the Stdout of the `git log ...` command
	scanner       *bufio.Scanner
	stderr        io.ReadCloser
	currentCommit *Commit
}

func (i *commitIterator) readUntilCompleteCommit() (*Commit, error) {
	var inCommitBody bool
	for i.scanner.Scan() {
		line := i.scanner.Text()

		switch {
		case strings.HasPrefix(line, commitHashPrefix):
			commitToReturn := i.currentCommit
			i.currentCommit = &Commit{
				SHA:       strings.TrimPrefix(line, commitHashPrefix),
				Author:    Event{},
				Committer: Event{},
				Parents:   []string{},
				Stats:     make([]Stat, 0),
			}

			if commitToReturn != nil { // if we're seeing a new commit but already have a current commit, we've finished a commit
				return commitToReturn, nil
			}
		case strings.HasPrefix(line, treeHashPrefix):
			i.currentCommit.Tree = strings.TrimPrefix(line, treeHashPrefix)
		case strings.HasPrefix(line, parentHashesPrefix):
			i.currentCommit.Parents = append(i.currentCommit.Parents, strings.TrimPrefix(line, parentHashesPrefix))
		case strings.HasPrefix(line, authorNamePrefix):
			i.currentCommit.Author.Name = strings.TrimPrefix(line, authorNamePrefix)
		case strings.HasPrefix(line, authorEmailPrefix):
			i.currentCommit.Author.Email = strings.TrimPrefix(line, authorEmailPrefix)
		case strings.HasPrefix(line, authorDatePrefix):
			s := strings.TrimPrefix(line, authorDatePrefix)
			if t, err := time.Parse(isoDataFmtStr, s); err != nil {
				return nil, err
			} else {
				i.currentCommit.Author.When = t
			}
		case strings.HasPrefix(line, committerNamePrefix):
			i.currentCommit.Committer.Name = strings.TrimPrefix(line, committerNamePrefix)
		case strings.HasPrefix(line, committerEmailPrefix):
			i.currentCommit.Committer.Email = strings.TrimPrefix(line, committerEmailPrefix)
		case strings.HasPrefix(line, committerDatePrefix):
			s := strings.TrimPrefix(line, committerDatePrefix)
			if t, err := time.Parse(isoDataFmtStr, s); err != nil {
				return nil, err
			} else {
				i.currentCommit.Committer.When = t
			}
		case strings.HasPrefix(line, commitBodyPrefix):
			inCommitBody = true
			s := strings.TrimPrefix(line, commitBodyPrefix)
			i.currentCommit.Message = s + "\n"
		case strings.HasPrefix(line, string([]byte{0})):
			inCommitBody = false
		default:
			if inCommitBody {
				i.currentCommit.Message += line + "\n"
				continue
			} else {
				i.currentCommit.hasTrailingNewline = true
			}
			s := strings.Split(line, "\t")
			if len(s) != 3 {
				continue
			}
			var additions int
			var deletions int
			var err error
			if s[0] != "-" {
				additions, err = strconv.Atoi(s[0])
				if err != nil {
					return nil, err
				}
			} else {
				additions = -1
			}
			if s[1] != "-" {
				deletions, err = strconv.Atoi(s[1])
				if err != nil {
					return nil, err
				}
			} else {
				deletions = -1
			}
			i.currentCommit.Stats = append(i.currentCommit.Stats, Stat{
				FilePath:  s[2],
				Additions: additions,
				Deletions: deletions,
			})
		}
	}

	return i.currentCommit, io.EOF
}

// Next moves the iterator and returns the next *Commit (or error)
// If the returned *Commit and err are nil, then iteration is complete.
func (i *commitIterator) Next() (*Commit, error) {
	if commit, err := i.readUntilCompleteCommit(); err != nil {
		if errors.Is(err, io.EOF) {
			i.currentCommit = nil
			return commit, nil
		}
		return nil, err
	} else {
		return commit, nil
	}
}

// Exec runs the git log command against a repository on disk and returns an iterator
// for walking over all the commits returned.
func Exec(ctx context.Context, repoPath string, options ...Option) (*commitIterator, error) {
	o := &execOptions{}
	for _, option := range options {
		option(o)
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("could not find git: %w", err)
	}

	args := []string{"log"}
	args = append(args, fmt.Sprintf("--format=%s", buildFormatString()), "--no-decorate", "-w")
	if o.NoMerges {
		args = append(args, "--no-merges")
	}

	if o.FirstParent {
		args = append(args, "--first-parent")
	}

	if o.Order != "" {
		args = append(args, string(o.Order))
	}

	if o.M {
		args = append(args, "-m")
	}

	if o.FileFilter != "" {
		args = append(args, o.FileFilter)
	}

	if o.Stats {
		args = append(args, "--numstat")
	}

	cmd := exec.CommandContext(ctx, gitPath, args...)
	cmd.Dir = repoPath

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)

	// this is a custom split function based off the default bufio.ScanLines one
	// https://cs.opensource.google/go/go/+/refs/tags/go1.19.1:src/bufio/scan.go;l=350;drc=18888751828c329ddf5efdd7ec1b39adf0b6ea00
	// we need to do this because we want to preserve the newlines in commit messages
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	})

	iter := &commitIterator{
		scanner: scanner,
		stderr:  stderr,
	}

	return iter, nil
}

// String returns a string representation of a commit that looks like the output from the original git log command
// where --format=... uses the custom format template we define in buildFormatString().
func (c *Commit) String() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%s%s\n", commitHashPrefix, c.SHA))
	s.WriteString(fmt.Sprintf("%s%s\n", treeHashPrefix, c.Tree))
	s.WriteString(fmt.Sprintf("%s%s\n", parentHashesPrefix, strings.Join(c.Parents, " ")))

	s.WriteString(fmt.Sprintf("%s%s\n", authorNamePrefix, c.Author.Name))
	s.WriteString(fmt.Sprintf("%s%s\n", authorEmailPrefix, c.Author.Email))
	s.WriteString(fmt.Sprintf("%s%s\n", authorDatePrefix, c.Author.When.Format(isoDataFmtStr)))

	s.WriteString(fmt.Sprintf("%s%s\n", committerNamePrefix, c.Committer.Name))
	s.WriteString(fmt.Sprintf("%s%s\n", committerEmailPrefix, c.Committer.Email))
	s.WriteString(fmt.Sprintf("%s%s\n", committerDatePrefix, c.Committer.When.Format(isoDataFmtStr)))

	message := c.Message
	// message = strings.TrimSuffix(message, "\n")

	s.WriteString(fmt.Sprintf("%s%s\x00\n", commitBodyPrefix, message))

	if len(c.Stats) > 0 || c.hasTrailingNewline {
		s.WriteString("\n")
	}

	for _, stat := range c.Stats {
		additions := fmt.Sprintf("%d", stat.Additions)
		deletions := fmt.Sprintf("%d", stat.Deletions)
		if additions == "-1" {
			additions = "-"
		}
		if deletions == "-1" {
			deletions = "-"
		}
		s.WriteString(fmt.Sprintf("%s\t%s\t%s\n", additions, deletions, stat.FilePath))
	}

	return s.String()
}
