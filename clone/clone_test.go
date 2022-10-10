package clone

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func TestSingleArgsOK(t *testing.T) {
	type test struct {
		options      []Option
		flags        []string
		path         string
		spectedError bool
	}

	tests := []test{
		{options: []Option{WithRejectShallow(true)}, flags: []string{"--reject-shallow"}, spectedError: false},
		{options: []Option{WithNoRejectShallow(true)}, flags: []string{"--no-reject-shallow"}, spectedError: false},
		{options: []Option{WithNoCheckout(true)}, flags: []string{"--no-checkout"}, spectedError: false},
		{options: []Option{WithBare(true)}, flags: []string{"--bare"}, spectedError: false},
		{options: []Option{WithMirror(true)}, flags: []string{"--mirror"}, spectedError: false},
		{options: []Option{WithLocal(true)}, flags: []string{"--local"}, spectedError: false},
		{options: []Option{WithNoHardLinks(true)}, flags: []string{"--no-hardlinks"}, spectedError: false},
		{options: []Option{WithShared(true)}, flags: []string{"--shared"}, spectedError: false},
		{options: []Option{WithDissociate(true)}, flags: []string{"--dissociate"}, spectedError: false},
		{options: []Option{WithQuiet(true)}, flags: []string{"--quiet"}, spectedError: false},
		{options: []Option{WithVerbose(true)}, flags: []string{"--verbose"}, spectedError: false},
		{options: []Option{WithProgress(true)}, flags: []string{"--progress"}, spectedError: false},
		{options: []Option{WithSparse(true)}, flags: []string{"--sparse"}, spectedError: false},
		{options: []Option{WithSingleBranch(true)}, flags: []string{"--single-branch"}, spectedError: false},
		{options: []Option{WithNoSingleBranch(true)}, flags: []string{"--no-single-branch"}, spectedError: false},
		{options: []Option{WithNoTags(true)}, flags: []string{"--no-tags"}, spectedError: false},
		{options: []Option{WithNoShallowSubmodules(true)}, flags: []string{"--no-shallow-submodules"}, spectedError: false},
		{options: []Option{WithRemoteSubmodules(true)}, flags: []string{"--remote-submodules"}, spectedError: false},
		{options: []Option{WithNoRemoteSubmodules(true)}, flags: []string{"--no-remote-submodules"}, spectedError: false},
		{options: []Option{WithFilter("some-filter")}, flags: []string{"--filter=some-filter"}, spectedError: false},
		{options: []Option{WithRecurseSubmodules("some-string")}, flags: []string{"--recurse-submodules=some-string"}, spectedError: false},
		{options: []Option{WithAlsoFilterSubmodules(true)}, flags: []string{"--also-filter-submodules"}, spectedError: false},
		{options: []Option{WithServerOptions([]string{"a", "b"})}, flags: []string{"--server-option", "a", "--server-option", "b"}, spectedError: false},
		{options: []Option{WithShallowExclude([]string{"a", "b"})}, flags: []string{"--shallow-exclude=a", "--shallow-exclude=b"}, spectedError: false},
		{options: []Option{WithSeparateGitDir("some-string")}, flags: []string{"--separate-git-dir=some-string"}, spectedError: false},
		{options: []Option{WithShallowSince("some-string")}, flags: []string{"--shallow-since", "some-string"}, spectedError: false},
		{options: []Option{WithTemplate("some-string")}, flags: []string{"--template", "some-string"}, spectedError: false},
		{options: []Option{WithUploadPack("some-string")}, flags: []string{"--upload-pack", "some-string"}, spectedError: false},
		{options: []Option{WithOrigin("some-string")}, flags: []string{"--origin", "some-string"}, spectedError: false},
		{options: []Option{WithReference("some-string")}, flags: []string{"--reference", "some-string"}, spectedError: false},
		{options: []Option{WithReferenceIfAble("some-string")}, flags: []string{"--reference-if-able", "some-string"}, spectedError: false},
		{options: []Option{WithBranch("some-string")}, flags: []string{"--branch", "some-string"}, spectedError: false},
		{options: []Option{WithDepth(5)}, flags: []string{"--depth=5"}, spectedError: false},
		{options: []Option{WithJobs(5)}, flags: []string{"--jobs=5"}, spectedError: false},
		{options: []Option{WithConfig([]ConfigKV{{Key: "a", Value: "b"}, {Key: "c", Value: "d"}})}, flags: []string{"--config=a=b", "--config=c=d"}, spectedError: false},
		{options: []Option{WithBare(true)}, flags: []string{"--bare"}, path: "***mergestat*repo*3642228556", spectedError: true},
	}

	for _, tc := range tests {
		t.Run(strings.Join(tc.flags, ","), func(t *testing.T) {

			ctx := context.Background()
			o := &execOptions{}

			// apply the options
			for _, opt := range tc.options {
				opt(o)
			}

			got := flagArgsFromOptions(o)

			if !reflect.DeepEqual(got, tc.flags) {
				t.Errorf("got %v, want %v", got, tc.flags)
			}

			if tc.spectedError {

				err := Exec(ctx, "https://github.com/mergestat/gitutils", tc.path, tc.options[0])

				if err == nil {
					t.Fatal(err)
				}
			}
		})
	}
}
