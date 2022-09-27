package clone

import (
	"context"
	"fmt"
	"os/exec"
)

type execOptions struct {
	RejectShallow        bool
	NoRejectShallow      bool
	NoCheckout           bool
	Bare                 bool
	Mirror               bool
	Local                bool
	NoHardLinks          bool
	Shared               bool
	Dissociate           bool
	Quiet                bool
	Verbose              bool
	Progress             bool
	Sparse               bool
	AlsoFilterSubModules bool
	NoSingleBranch       bool
	SingleBranch         bool
	NoTags               bool
	ShallowSubmodules    bool
	NoShallowSubmodules  bool
	RemoteSubmodules     bool
	NoRemoteSubmodules   bool
	UploadPack           string
	Origin               string
	Branch               string
	Reference            string
	ReferenceIfAble      string
	Filter               string
	Template             string
	ShallowSince         string
	RecursiveSubmodules  string
	SeparateGitDir       string
	ShallowExclude       []string
	ServerOptions        []string
	Depth                int
	Jobs                 int
	Config               map[string]string
}

type Option func(o *execOptions)

// WithRejectShallow sets the --reject-shallow flag
func WithRejectShallow(rejectShallow bool) Option {
	return func(o *execOptions) {
		o.RejectShallow = rejectShallow
	}
}

// WithNoRejectShallow sets the --no-reject-shallow flag
func WithNoRejectShallow(noRejectShallow bool) Option {
	return func(o *execOptions) {
		o.NoRejectShallow = noRejectShallow
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

// WithReference sets the --reference <repository> flag
func WithReference(reference string) Option {
	return func(o *execOptions) {
		o.Reference = reference
	}
}

// WithReferenceIfAble sets the --reference-if-able <repository> flag
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

// WithSparse sets the --sparse flag
func WithSparse(sparse bool) Option {
	return func(o *execOptions) {
		o.Sparse = sparse
	}
}

// Withfilter sets the --filter <filter> flag
func WithFilter(filter string) Option {
	return func(o *execOptions) {
		o.Filter = filter
	}
}

// WithOrigin sets the --origin <name> flag
func WithOrigin(origin string) Option {
	return func(o *execOptions) {
		o.Origin = origin
	}
}

// WithUploadPack sets the --upload-pack <dir> flag
func WithUploadPack(uploadPack string) Option {
	return func(o *execOptions) {
		o.UploadPack = uploadPack
	}
}

// WithTemplate sets the --template <template-dir> flag
func WithTemplate(template string) Option {
	return func(o *execOptions) {
		o.Template = template
	}
}

// WithDepth sets the --depth <depth> flag
func WithDepth(depth int) Option {
	return func(o *execOptions) {
		o.Depth = depth
	}
}

// WithAlsoFilterSubmodules sets the --also-filter-submodules flag .For usage Refs :https://git-scm.com/docs/git-clone
func WithAlsoFilterSubmodules(alsoFilterSubModules bool) Option {
	return func(o *execOptions) {
		o.AlsoFilterSubModules = alsoFilterSubModules
	}
}

// WithShallowSince sets the --shallow-since <date> flag
func WithShallowSince(shallowSince string) Option {
	return func(o *execOptions) {
		o.ShallowSince = shallowSince
	}
}

// WithNoSingleBranch sets the --no-single-branch flag
func WithNoSingleBranch(noSingleBranch bool) Option {
	return func(o *execOptions) {
		o.NoSingleBranch = noSingleBranch
	}
}

// WithSingleBranch sets the --single-branch flag
func WithSingleBranch(singleBranch bool) Option {
	return func(o *execOptions) {
		o.SingleBranch = singleBranch
	}
}

// WithNoTags sets the --no-tags flag
func WithNoTags(noTags bool) Option {
	return func(o *execOptions) {
		o.NoTags = noTags
	}
}

// WithRecurseSubmodules sets the --recurse-modules <pathspec> flag
func WithRecurseSubmodules(recurseSubmodules string) Option {
	return func(o *execOptions) {
		o.RecursiveSubmodules = recurseSubmodules
	}
}

// WithNoShallowSubmodules sets the --no-shallow-submodules flag
func WithNoShallowSubmodules(noShallowSubmodules bool) Option {
	return func(o *execOptions) {
		o.NoShallowSubmodules = noShallowSubmodules
	}
}

// WithShallowSubmodules sets the --shallow-submodules flag
func WithShallowSubmodules(shallowSubmodules bool) Option {
	return func(o *execOptions) {
		o.ShallowSubmodules = shallowSubmodules
	}
}

// WithRemoteSubmodules sets the --remote-submodules flag
func WithRemoteSubmodules(remoteSubmodules bool) Option {
	return func(o *execOptions) {
		o.RemoteSubmodules = remoteSubmodules
	}
}

// WithNoRemoteSubmodules sets the --no-remote-submodules flag
func WithNoRemoteSubmodules(noRemoteSubmodules bool) Option {
	return func(o *execOptions) {
		o.NoRemoteSubmodules = noRemoteSubmodules
	}
}

// WithSeparateGitDir sets the --separate-git-dir <dir> flag
func WithSeparateGitDir(separateGitDir string) Option {
	return func(o *execOptions) {
		o.SeparateGitDir = separateGitDir
	}
}

// WithConfig sets the --config <key>=<value> flag
func WithConfig(config map[string]string) Option {
	return func(o *execOptions) {
		o.Config = config
	}
}

// WithJobs sets the --jobs <jobs> flag
func WithJobs(jobs int) Option {
	return func(o *execOptions) {
		o.Jobs = jobs
	}
}

// WithShallowExclude sets the --shallow-exclude <revision> flag
func WithShallowExclude(shallowExclude []string) Option {
	return func(o *execOptions) {
		o.ShallowExclude = shallowExclude
	}
}

// WithServerOptions sets the --server-option <option> flag
func WithServerOptions(ServerOptions []string) Option {
	return func(o *execOptions) {
		o.ServerOptions = ServerOptions
	}
}

// flagArgsFromOptions returns a slice of flags from the given options struct
func flagArgsFromOptions(o *execOptions) []string {
	var args []string

	if o.RejectShallow {
		args = append(args, "--reject-shallow")
	}

	if o.NoRejectShallow {
		args = append(args, "--no-reject-shallow")
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

	if o.Sparse {
		args = append(args, "--sparse")
	}

	if o.SingleBranch {
		args = append(args, "--single-branch")
	}

	if o.NoSingleBranch {
		args = append(args, "--no-single-branch")
	}

	if o.NoTags {
		args = append(args, "--no-tags")
	}

	if o.ShallowSubmodules {
		args = append(args, "--shallow-submodules")
	}

	if o.NoShallowSubmodules {
		args = append(args, "--no-shallow-submodules")
	}

	if o.RemoteSubmodules {
		args = append(args, "--remote-submodules")
	}

	if o.NoRemoteSubmodules {
		args = append(args, "--no-remote-submodules")
	}

	if len(o.Filter) > 0 {
		args = append(args, fmt.Sprintf("--filter=%s", o.Filter))
	}

	if len(o.RecursiveSubmodules) > 0 {
		args = append(args, fmt.Sprintf("--recurse-submodules=%s", o.RecursiveSubmodules))
	}

	if o.AlsoFilterSubModules {
		args = append(args, "--also-filter-submodules")
	}

	if len(o.ServerOptions) > 0 {
		for _, option := range o.ServerOptions {
			args = append(args, "--server-option", option)
		}
	}

	if len(o.ShallowExclude) > 0 {
		for _, option := range o.ShallowExclude {
			args = append(args, fmt.Sprintf("--shallow-exclude=%s", option))
		}
	}

	if len(o.SeparateGitDir) > 0 {
		args = append(args, fmt.Sprintf("--separate-git-dir=%s", o.SeparateGitDir))
	}

	if len(o.ShallowSince) > 0 {
		args = append(args, "--shallow-since", o.ShallowSince)
	}

	if len(o.Template) > 0 {
		args = append(args, "--template", o.Template)
	}

	if len(o.UploadPack) > 0 {
		args = append(args, "--upload-pack", o.UploadPack)
	}

	if len(o.Origin) > 0 {
		args = append(args, "--origin", o.Origin)
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

	if o.Depth > 0 {
		args = append(args, fmt.Sprintf("--depth=%d", o.Depth))
	}

	if o.Jobs > 0 {
		args = append(args, fmt.Sprintf("--jobs=%d", o.Jobs))
	}

	if len(o.Config) > 0 {
		for k, v := range o.Config {
			config := fmt.Sprintf("--config=%s=%s", k, v)
			args = append(args, config)
		}
	}

	return args
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
	args = append(args, flagArgsFromOptions(o)...)

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
