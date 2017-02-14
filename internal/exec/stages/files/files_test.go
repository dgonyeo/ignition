// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package files

import (
	"sort"
	"testing"

	"github.com/kylelemons/godebug/pretty"

	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/exec/util"
	"github.com/coreos/ignition/internal/log"
)

func TestMapEntriesToFilesystems(t *testing.T) {
	type in struct {
		config types.Config
	}
	type out struct {
		files map[types.Filesystem][]filesystemEntry
		err   error
	}

	fs1 := "/fs1"
	fs2 := "/fs2"

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{config: types.Config{}},
			out: out{files: map[types.Filesystem][]filesystemEntry{}},
		},
		{
			in:  in{config: types.Config{Storage: types.Storage{Files: []types.File{{Node: types.Node{Filesystem: "foo"}}}}}},
			out: out{err: ErrFilesystemUndefined},
		},
		{
			in: in{config: types.Config{Storage: types.Storage{
				Filesystems: []types.Filesystem{{Name: "fs1"}},
				Files: []types.File{
					{Node: types.Node{Filesystem: "fs1", Path: "/foo"}},
					{Node: types.Node{Filesystem: "fs1", Path: "/bar"}},
				},
			}}},
			out: out{files: map[types.Filesystem][]filesystemEntry{types.Filesystem{Name: "fs1"}: {
				fileEntry(types.File{Node: types.Node{Filesystem: "fs1", Path: "/foo"}}),
				fileEntry(types.File{Node: types.Node{Filesystem: "fs1", Path: "/bar"}}),
			}}},
		},
		{
			in: in{config: types.Config{Storage: types.Storage{
				Filesystems: []types.Filesystem{{Name: "fs1", Path: &fs1}, {Name: "fs2", Path: &fs2}},
				Files: []types.File{
					{Node: types.Node{Filesystem: "fs1", Path: "/foo"}},
					{Node: types.Node{Filesystem: "fs2", Path: "/bar"}},
				},
			}}},
			out: out{files: map[types.Filesystem][]filesystemEntry{
				types.Filesystem{Name: "fs1", Path: &fs1}: {fileEntry(types.File{Node: types.Node{Filesystem: "fs1", Path: "/foo"}})},
				types.Filesystem{Name: "fs2", Path: &fs2}: {fileEntry(types.File{Node: types.Node{Filesystem: "fs2", Path: "/bar"}})},
			}},
		},
		{
			in: in{config: types.Config{Storage: types.Storage{
				Filesystems: []types.Filesystem{{Name: "fs1"}, {Name: "fs1", Path: &fs1}},
				Files: []types.File{
					{Node: types.Node{Filesystem: "fs1", Path: "/foo"}},
					{Node: types.Node{Filesystem: "fs1", Path: "/bar"}},
				},
			}}},
			out: out{files: map[types.Filesystem][]filesystemEntry{
				types.Filesystem{Name: "fs1", Path: &fs1}: {
					fileEntry(types.File{Node: types.Node{Filesystem: "fs1", Path: "/foo"}}),
					fileEntry(types.File{Node: types.Node{Filesystem: "fs1", Path: "/bar"}}),
				},
			}},
		},
	}

	for i, test := range tests {
		logger := log.New()
		files, err := stage{Util: util.Util{Logger: &logger}}.mapEntriesToFilesystems(test.in.config)
		if str := pretty.Compare(err, test.out.err); str != "" {
			t.Errorf("#%d: bad error:\n%s", i, str)
		}
		if str := pretty.Compare(files, test.out.files); str != "" {
			t.Errorf("#%d: bad map:\n%s", i, str)
		}
	}
}

func TestDirectorySort(t *testing.T) {
	type in struct {
		data []string
	}

	type out struct {
		data []string
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{data: []string{"/a/b/c/d/e/", "/a/b/c/d/", "/a/b/c/", "/a/b/", "/a/"}},
			out: out{data: []string{"/a/", "/a/b/", "/a/b/c/", "/a/b/c/d/", "/a/b/c/d/e/"}},
		},
		{
			in:  in{data: []string{"/a////b/c/d/e/", "/", "/a/b/c//d/", "/a/b/c/", "/a/b/", "/a/"}},
			out: out{data: []string{"/", "/a/", "/a/b/", "/a/b/c/", "/a/b/c//d/", "/a////b/c/d/e/"}},
		},
		{
			in:  in{data: []string{"/a/", "/a/../a/b", "/"}},
			out: out{data: []string{"/", "/a/", "/a/../a/b"}},
		},
	}

	for i, test := range tests {
		dirs := make([]types.Node, len(test.in.data))
		for j := range dirs {
			dirs[j].Path = test.in.data[j]
		}
		sort.Sort(ByDirectorySegments(dirs))
		outpaths := make([]string, len(test.in.data))
		for j, dir := range dirs {
			outpaths[j] = dir.Path
		}
		if str := pretty.Compare(outpaths, test.out.data); str != "" {
			t.Errorf("#%d: bad error:\n%s", i, str)
		}
	}
}
