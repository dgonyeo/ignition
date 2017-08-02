// Copyright 2017 CoreOS, Inc.
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
	"github.com/coreos/ignition/tests/types"
)

func Create_hard_link_on_root() types.Test {
	var mntDevices []types.MntDevice

	name := "Create a Hard Link on the Root Filesystem"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices = nil
	config := `{
	  "ignition": { "version": "2.1.0" },
	  "storage": {
	    "files": [{
	      "filesystem": "root",
	      "path": "/foo/target",
	      "contents": {
	        "source": "http://127.0.0.1:8080/contents"
	      }
	    }],
	    "links": [{
	      "filesystem": "root",
	      "path": "/foo/bar",
		  "target": "/foo/target",
		  "hard": true
	    }]
	  }
	}`
	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Path: "foo",
				Name: "target",
			},
			Contents: []string{"asdf\nfdsa"},
		},
	})
	out[0].Partitions.AddLinks("ROOT", []types.Link{
		{
			Node: types.Node{
				Path: "foo",
				Name: "bar",
			},
			Target: "/foo/target",
			Hard:   true,
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}

func Create_symlink_on_root() types.Test {
	var mntDevices []types.MntDevice

	name := "Create a Symlink on the Root Filesystem"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices = nil
	config := `{
	  "ignition": { "version": "2.1.0" },
	  "storage": {
	    "links": [{
	      "filesystem": "root",
	      "path": "/foo/bar",
		  "target": "/foo/target",
		  "hard": false
	    }]
	  }
	}`
	out[0].Partitions.AddLinks("ROOT", []types.Link{
		{
			Node: types.Node{
				Name: "bar",
				Path: "foo",
			},
			Target: "/foo/target",
			Hard:   false,
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}
