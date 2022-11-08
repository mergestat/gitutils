package blame

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Blame represents the blame of a particular line
type Blame struct {
	SHA            string
	OriginalLineNo int
	FinalLineNo    int
	LinesInGroup   int
	Author         *Event
	Committer      *Event
	Line           string
	Summary        string
	Boundary       bool
	Previous       string // TODO(patrickdevivo) split the SHA and filename out
	Filename       string
}

// Event represents the who and when of a commit event
type Event struct {
	Name  string
	Email string
	When  time.Time
}

func (blame *Blame) String() string {
	return fmt.Sprintf("%s: %s <%s>", blame.SHA, blame.Author.Name, blame.Author.Email)
}

func (event *Event) String() string {
	return fmt.Sprintf("%s <%s>", event.Name, event.Email)
}

// Result is a mapping of line numbers to blames for a given file
type Result []*Blame

func parseLinePorcelain(reader io.Reader, o *execOptions) (Result, error) {
	scanner := bufio.NewScanner(reader)

	if o.ScannerBuffer != nil {
		scanner.Buffer(o.ScannerBuffer, o.ScannerBufferMax)
	}

	res := make(Result, 0)

	const (
		author     = "author "
		authorMail = "author-mail "
		authorTime = "author-time "
		authorTZ   = "author-tz "

		committer     = "committer "
		committerMail = "committer-mail "
		committerTime = "committer-time "
		committerTZ   = "committer-tz "

		summary    = "summary "
		boundary   = "boundary"
		previous   = "previous "
		filename   = "filename "
		linePrefix = "\t"
	)

	var currentBlame *Blame
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, author):
			currentBlame.Author.Name = strings.TrimPrefix(line, author)
		case strings.HasPrefix(line, authorMail):
			s := strings.TrimPrefix(line, authorMail)
			currentBlame.Author.Email = strings.Trim(s, "<>")
		case strings.HasPrefix(line, authorTime):
			timeString := strings.TrimPrefix(line, authorTime)
			i, err := strconv.ParseInt(timeString, 10, 64)
			if err != nil {
				return nil, err
			}
			currentBlame.Author.When = time.Unix(i, 0)
		case strings.HasPrefix(line, authorTZ):
			tzString := strings.TrimPrefix(line, authorTZ)
			parsed, err := time.Parse("-0700", tzString)
			if err != nil {
				return nil, err
			}
			loc := parsed.Location()
			currentBlame.Author.When = currentBlame.Author.When.In(loc)
		case strings.HasPrefix(line, committer):
			currentBlame.Committer.Name = strings.TrimPrefix(line, committer)
		case strings.HasPrefix(line, committerMail):
			s := strings.TrimPrefix(line, committerMail)
			currentBlame.Committer.Email = strings.Trim(s, "<>")
		case strings.HasPrefix(line, committerTime):
			timeString := strings.TrimPrefix(line, committerTime)
			i, err := strconv.ParseInt(timeString, 10, 64)
			if err != nil {
				return nil, err
			}
			currentBlame.Committer.When = time.Unix(i, 0)
		case strings.HasPrefix(line, committerTZ):
			tzString := strings.TrimPrefix(line, committerTZ)
			parsed, err := time.Parse("-0700", tzString)
			if err != nil {
				return nil, err
			}
			loc := parsed.Location()
			currentBlame.Committer.When = currentBlame.Committer.When.In(loc)
		case strings.HasPrefix(line, summary):
			currentBlame.Summary = strings.TrimPrefix(line, summary)
		case strings.HasPrefix(line, boundary):
			currentBlame.Boundary = true
		case strings.HasPrefix(line, previous):
			currentBlame.Previous = strings.TrimPrefix(line, previous)
		case strings.HasPrefix(line, filename):
			currentBlame.Filename = strings.TrimPrefix(line, filename)
		case strings.HasPrefix(line, linePrefix):
			currentBlame.Line = strings.TrimPrefix(line, linePrefix)
		case len(strings.Split(line, " ")[0]) == 40: // if the first string sep by a space is 40 chars long, it's probably the commit header
			// there's an existing currentBlame, add it to the response
			if currentBlame != nil {
				res = append(res, currentBlame)
			}

			// reset the currentBlame to an empty struct, to be filled
			currentBlame = &Blame{
				Author:    &Event{},
				Committer: &Event{},
			}

			split := strings.Split(line, " ")
			sha := split[0]
			var err error

			if currentBlame.OriginalLineNo, err = strconv.Atoi(split[1]); err != nil {
				return nil, err
			}

			if currentBlame.FinalLineNo, err = strconv.Atoi(split[2]); err != nil {
				return nil, err
			}

			if len(split) > 3 {
				if currentBlame.LinesInGroup, err = strconv.Atoi(split[3]); err != nil {
					return nil, err
				}
			}

			// set the SHA since we have it now
			currentBlame.SHA = sha
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if currentBlame != nil {
		res = append(res, currentBlame)
	}

	return res, nil
}

type Option func(o *execOptions)

type execOptions struct {
	Revision         string
	ScannerBuffer    []byte
	ScannerBufferMax int
}

func WithRevision(revision string) Option {
	return func(o *execOptions) {
		o.Revision = revision
	}
}

// WithScannerBuffer sets the buffer and max buffer size for the scanner
// used when reading the output of the git blame command.
func WithScannerBuffer(buf []byte, max int) Option {
	return func(o *execOptions) {
		o.ScannerBuffer = buf
		o.ScannerBufferMax = max
	}
}

// Exec uses git to lookup the blame of a file, given the supplied options
func Exec(ctx context.Context, repoPath, filePath string, options ...Option) (Result, error) {
	o := &execOptions{}
	for _, option := range options {
		option(o)
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("could not find git: %w", err)
	}

	args := []string{"blame", "--line-porcelain", filePath}
	if o.Revision != "" {
		args = append(args, o.Revision)
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

	res, err := parseLinePorcelain(stdout, o)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		// TODO report the error message (on stderr) if it exists?
		// if exitErr, ok := err.(*exec.ExitError); ok {
		// 	string(exitErr.Stderr)
		// }
		return nil, err
	}

	return res, nil
}
