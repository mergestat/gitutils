package clone

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func TestSingleArgsOK(t *testing.T) {
	type test struct {
		options []Option
		flags   []string
	}

	tests := []test{
		{options: []Option{WithRejectShallow(true)}, flags: []string{"--reject-shallow"}},
		{options: []Option{WithNoRejectShallow(true)}, flags: []string{"--no-reject-shallow"}},
		{options: []Option{WithNoCheckout(true)}, flags: []string{"--no-checkout"}},
		{options: []Option{WithBare(true)}, flags: []string{"--bare"}},
		{options: []Option{WithMirror(true)}, flags: []string{"--mirror"}},
		{options: []Option{WithLocal(true)}, flags: []string{"--local"}},
		{options: []Option{WithNoHardLinks(true)}, flags: []string{"--no-hardlinks"}},
		{options: []Option{WithShared(true)}, flags: []string{"--shared"}},
		{options: []Option{WithDissociate(true)}, flags: []string{"--dissociate"}},
		{options: []Option{WithQuiet(true)}, flags: []string{"--quiet"}},
		{options: []Option{WithVerbose(true)}, flags: []string{"--verbose"}},
		{options: []Option{WithProgress(true)}, flags: []string{"--progress"}},
		{options: []Option{WithSparse(true)}, flags: []string{"--sparse"}},
		{options: []Option{WithSingleBranch(true)}, flags: []string{"--single-branch"}},
		{options: []Option{WithNoSingleBranch(true)}, flags: []string{"--no-single-branch"}},
		{options: []Option{WithNoTags(true)}, flags: []string{"--no-tags"}},
		{options: []Option{WithNoShallowSubmodules(true)}, flags: []string{"--no-shallow-submodules"}},
		{options: []Option{WithRemoteSubmodules(true)}, flags: []string{"--remote-submodules"}},
		{options: []Option{WithNoRemoteSubmodules(true)}, flags: []string{"--no-remote-submodules"}},
		{options: []Option{WithFilter("some-filter")}, flags: []string{"--filter=some-filter"}},
		{options: []Option{WithRecurseSubmodules("some-string")}, flags: []string{"--recurse-submodules=some-string"}},
		{options: []Option{WithAlsoFilterSubmodules(true)}, flags: []string{"--also-filter-submodules"}},
		{options: []Option{WithServerOptions([]string{"a", "b"})}, flags: []string{"--server-option", "a", "--server-option", "b"}},
		{options: []Option{WithShallowExclude([]string{"a", "b"})}, flags: []string{"--shallow-exclude=a", "--shallow-exclude=b"}},
		{options: []Option{WithSeparateGitDir("some-string")}, flags: []string{"--separate-git-dir=some-string"}},
		{options: []Option{WithShallowSince("some-string")}, flags: []string{"--shallow-since", "some-string"}},
		{options: []Option{WithTemplate("some-string")}, flags: []string{"--template", "some-string"}},
		{options: []Option{WithUploadPack("some-string")}, flags: []string{"--upload-pack", "some-string"}},
		{options: []Option{WithOrigin("some-string")}, flags: []string{"--origin", "some-string"}},
		{options: []Option{WithReference("some-string")}, flags: []string{"--reference", "some-string"}},
		{options: []Option{WithReferenceIfAble("some-string")}, flags: []string{"--reference-if-able", "some-string"}},
		{options: []Option{WithBranch("some-string")}, flags: []string{"--branch", "some-string"}},
		{options: []Option{WithDepth(5)}, flags: []string{"--depth=5"}},
		{options: []Option{WithJobs(5)}, flags: []string{"--jobs=5"}},
		{options: []Option{WithConfig([]ConfigKV{{Key: "a", Value: "b"}, {Key: "c", Value: "d"}})}, flags: []string{"--config=a=b", "--config=c=d"}},
	}

	for _, tc := range tests {
		t.Run(strings.Join(tc.flags, ","), func(t *testing.T) {

			o := &execOptions{}

			// apply the options
			for _, opt := range tc.options {
				opt(o)
			}

			got := flagArgsFromOptions(o)

			if !reflect.DeepEqual(got, tc.flags) {
				t.Errorf("got %v, want %v", got, tc.flags)
			}

		})
	}
}

func TestInvalidPath(t *testing.T) {
	type test struct {
		options      []Option
		flags        []string
		path         string
		spectedError bool
	}

	tests := []test{
		{options: []Option{WithBare(true)}, flags: []string{"--bare"}, path: "/tmp/mergestat-repo-3642228556", spectedError: true},
	}

	for _, tc := range tests {
		t.Run(strings.Join(tc.flags, ","), func(t *testing.T) {

			ctx := context.Background()

			err := Exec(ctx, "https://github.com/mergestat/gitutils", tc.path, tc.options[0])

			if err == nil {
				t.Fatal(err)
			}

		})
	}
}
