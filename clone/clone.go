package clone

import (
	"context"
	"fmt"
	"os/exec"
)

type execOptions struct {
	RejectShallow bool
	NoCheckout    bool
	Bare          bool
	Mirror        bool
	Branch        bool
}

type Option func(o *execOptions)

// WithRejectShallow sets the --reject-shallow flag
func WithRejectShallow(rejectShallow bool) Option {
	return func(o *execOptions) {
		o.RejectShallow = rejectShallow
	}
}

// WithNoCheckout sets the --no-checkout flag
func WithNoCheckout(noCheckout bool) Option {
	return func(o *execOptions) {
		o.NoCheckout = noCheckout
	}
}

// WithBare sets the --bare flag
func WithBare(bare bool) Option {
	return func(o *execOptions) {
		o.Bare = bare
	}
}

// WithMirror sets the --mirror flag
func WithMirror(mirror bool) Option {
	return func(o *execOptions) {
		o.Mirror = mirror
	}
}

// WithBrach sets the --branch <name> flag
func WithBranch(branch bool) Option {
	return func(o *execOptions) {
		o.Branch = branch
	}
}

// Exec runs `git clone` using the os/exec standard library package.
func Exec(ctx context.Context, repo, dir string, flagArg string, options ...Option) error {
	o := &execOptions{}
	for _, option := range options {
		option(o)
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("could not find git: %w", err)
	}

	args := []string{"clone", repo, dir}

	if o.RejectShallow {
		args = append(args, "--reject-shallow")
	}

	if o.NoCheckout {
		args = append(args, "--no-checkout")
	}

	if o.Bare {
		args = append(args, "--bare")
	}

	if o.Mirror {
		args = append(args, "--mirror")
	}

	if o.Branch {
		args = append(args, "--branch", flagArg)
	}

	cmd := exec.CommandContext(ctx, gitPath, args...)

	_, err = cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
