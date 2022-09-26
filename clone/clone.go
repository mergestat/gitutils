package clone

import (
	"context"
	"fmt"
	"os/exec"
)

type execOptions struct {
	RejectShallow   bool
	NoCheckout      bool
	Bare            bool
	Mirror          bool
	Local           bool
	NoHardLinks     bool
	Shared          bool
	Dissociate      bool
	Quiet           bool
	Verbose         bool
	Progress        bool
	Branch          string
	Reference       string
	ReferenceIfAble string
	ServerOption    []string
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

// WithBranch sets the --branch <name> flag
func WithBranch(branch string) Option {
	return func(o *execOptions) {
		o.Branch = branch
	}
}

// WithLocal sets the --local flag
func WithLocal(local bool) Option {
	return func(o *execOptions) {
		o.Local = local
	}
}

// WithNoHardLinks sets the --no-hardlinks flag
func WithNoHardLinks(noHardLinks bool) Option {
	return func(o *execOptions) {
		o.NoHardLinks = noHardLinks
	}
}

// WithShared sets the --shared flag
func WithShared(shared bool) Option {
	return func(o *execOptions) {
		o.Shared = shared
	}
}

// WithReference sets the --reference flag
func WithReference(reference string) Option {
	return func(o *execOptions) {
		o.Reference = reference
	}
}

// WithReferenceIfAble sets the --reference-if-able flag
func WithReferenceIfAble(referenceIfAble string) Option {
	return func(o *execOptions) {
		o.ReferenceIfAble = referenceIfAble
	}
}

// WithDissociate sets the --dissociate flag
func WithDissociate(dissociate bool) Option {
	return func(o *execOptions) {
		o.Dissociate = dissociate
	}
}

// WithQuiet sets the --quiet flag
func WithQuiet(quiet bool) Option {
	return func(o *execOptions) {
		o.Quiet = quiet
	}
}

// WithVerbose sets the --verbose flag
func WithVerbose(verbose bool) Option {
	return func(o *execOptions) {
		o.Verbose = verbose
	}
}

// WithProgress sets the --progress flag
func WithProgress(progress bool) Option {
	return func(o *execOptions) {
		o.Progress = progress
	}
}

// WithServerOption sets the --server-option flag
func WithServerOption(ServerOption []string) Option {
	return func(o *execOptions) {
		o.ServerOption = ServerOption
	}
}

// Exec runs `git clone` using the os/exec standard library package.
func Exec(ctx context.Context, repo, dir string, options ...Option) error {
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

	if o.Local {
		args = append(args, "--local")
	}

	if o.NoHardLinks {
		args = append(args, "--no-hardlinks")
	}

	if o.Shared {
		args = append(args, "--shared")
	}

	if o.Dissociate {
		args = append(args, "--dissociate")
	}

	if o.Quiet {
		args = append(args, "--quiet")
	}

	if o.Verbose {
		args = append(args, "--verbose")
	}

	if o.Progress {
		args = append(args, "--progress")
	}

	if len(o.ServerOption) > 0 {
		for _, option := range o.ServerOption {
			args = append(args, "--server-option", option)
		}
	}

	if len(o.Reference) > 0 {
		args = append(args, "--reference", o.Reference)
	}

	if len(o.ReferenceIfAble) > 0 {
		args = append(args, "--reference-if-able", o.ReferenceIfAble)
	}

	if len(o.Branch) > 0 {
		args = append(args, "--branch", o.Branch)
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
