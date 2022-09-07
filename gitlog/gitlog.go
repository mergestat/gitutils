// Package gitlog ...
package gitlog

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// line prefixes for the `raw` formatted output
const (
	commitPrefix    = "commit "
	treePrefix      = "tree "
	parentPrefix    = "parent "
	authorPrefix    = "author "
	committerPrefix = "committer "
	gpgsigPrefix    = "gpgsig "
	messagePrefix   = "    "
)

// Commit represents a parsed commit from git log
type Commit struct {
	SHA       string
	Tree      string
	Parents   []string
	Author    Event
	Committer Event
	GPGSig    string
	Message   string
	Stats     []Stat
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

type commitIterator struct {
	// scanner is a Scanner produced from the Stdout of the `git log ...` command
	scanner       *bufio.Scanner
	stderr        io.ReadCloser
	currentCommit *Commit
}

func (i *commitIterator) readUntilCompleteCommit() (*Commit, error) {
	var inCommitMessage bool
	for i.scanner.Scan() {
		line := i.scanner.Text()
		switch {
		case strings.HasPrefix(line, commitPrefix):
			commitToReturn := i.currentCommit
			i.currentCommit = &Commit{
				SHA:       strings.TrimPrefix(line, commitPrefix),
				Author:    Event{},
				Committer: Event{},
				Parents:   []string{},
				Stats:     make([]Stat, 0),
			}

			if commitToReturn != nil { // if we're seeing a new commit but already have a current commit, we've finished a commit
				commitToReturn.Message = strings.TrimLeft(commitToReturn.Message, "\n")
				commitToReturn.Message = strings.TrimSuffix(commitToReturn.Message, "\n\n")
				return commitToReturn, nil
			}
		case strings.HasPrefix(line, treePrefix):
			i.currentCommit.Tree = strings.TrimPrefix(line, treePrefix)
		case strings.HasPrefix(line, parentPrefix):
			i.currentCommit.Parents = append(i.currentCommit.Parents, strings.TrimPrefix(line, parentPrefix))
		case strings.HasPrefix(line, authorPrefix):
			s := strings.TrimPrefix(line, authorPrefix)
			spl := strings.Split(s, " ")

			tz := spl[len(spl)-1]
			ts := spl[len(spl)-2]
			email := strings.Trim(spl[len(spl)-3], "<>")
			name := strings.Join(spl[:len(spl)-3], " ")

			i.currentCommit.Author.Email = email
			i.currentCommit.Author.Name = strings.TrimSpace(name)

			its, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				return nil, err
			}

			i.currentCommit.Author.When = time.Unix(its, 0)

			ptz, err := time.Parse("-0700", tz)
			if err != nil {
				return nil, err
			}
			i.currentCommit.Author.When = i.currentCommit.Author.When.In(ptz.Location())
		case strings.HasPrefix(line, committerPrefix):
			s := strings.TrimPrefix(line, committerPrefix)
			spl := strings.Split(s, " ")

			tz := spl[len(spl)-1]
			ts := spl[len(spl)-2]
			email := strings.Trim(spl[len(spl)-3], "<>")
			name := strings.Join(spl[:len(spl)-3], " ")

			i.currentCommit.Committer.Email = email
			i.currentCommit.Committer.Name = strings.TrimSpace(name)

			its, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				return nil, err
			}

			i.currentCommit.Committer.When = time.Unix(its, 0)

			ptz, err := time.Parse("-0700", tz)
			if err != nil {
				return nil, err
			}
			i.currentCommit.Committer.When = i.currentCommit.Committer.When.In(ptz.Location())
		case strings.HasPrefix(line, gpgsigPrefix):
			i.currentCommit.GPGSig += strings.TrimPrefix(line, gpgsigPrefix) + "\n"
		case strings.HasPrefix(line, messagePrefix):
			inCommitMessage = true
			i.currentCommit.Message += strings.TrimPrefix(line, messagePrefix) + "\n"
		case strings.HasPrefix(line, " "):
			i.currentCommit.GPGSig += strings.TrimPrefix(line, " ") + "\n"
		case line == "":
			if inCommitMessage {
				i.currentCommit.Message += "\n"
			}
		default:
			inCommitMessage = false
			s := strings.Split(line, "\t")
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

	if i.currentCommit != nil {
		i.currentCommit.Message = strings.TrimLeft(i.currentCommit.Message, "\n")
		i.currentCommit.Message = strings.TrimSuffix(i.currentCommit.Message, "\n\n")
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
	args = append(args, "--numstat", "--format=raw", "--no-decorate", "-w")
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

	iter := &commitIterator{
		scanner: bufio.NewScanner(stdout),
		stderr:  stderr,
	}

	return iter, nil
}

func prefixEveryLine(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("%s%s", prefix, line)
	}
	return strings.Join(lines, "\n")
}

// String returns a string representation of a commit that looks like the output of --format=raw
func (c *Commit) String() string {
	var output strings.Builder
	output.WriteString(fmt.Sprintf("commit %s\n", c.SHA))
	output.WriteString(fmt.Sprintf("tree %s\n", c.Tree))
	for _, parent := range c.Parents {
		output.WriteString(fmt.Sprintf("parent %s\n", parent))
	}
	output.WriteString(fmt.Sprintf("author %s <%s> %d %s\n", c.Author.Name, c.Author.Email, c.Author.When.Unix(), c.Author.When.Format("-0700")))
	output.WriteString(fmt.Sprintf("committer %s <%s> %d %s\n", c.Committer.Name, c.Committer.Email, c.Committer.When.Unix(), c.Committer.When.Format("-0700")))

	if c.GPGSig != "" {
		gpgsigLines := strings.Split(c.GPGSig, "\n")
		output.WriteString(fmt.Sprintf("gpgsig %s\n", gpgsigLines[0]))

		output.WriteString(prefixEveryLine(strings.Join(gpgsigLines[1:len(gpgsigLines)-1], "\n"), " "))
		output.WriteString("\n")
	}

	output.WriteString("\n")
	output.WriteString(strings.TrimSuffix(prefixEveryLine(c.Message, "    "), "    "))
	output.WriteString("\n")

	if len(c.Stats) > 0 {
		output.WriteString("\n")
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
		output.WriteString(fmt.Sprintf("%s\t%s\t%s\n", additions, deletions, stat.FilePath))
	}

	return output.String()
}
