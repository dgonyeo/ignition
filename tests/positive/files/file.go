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

func Create_file_on_root() types.Test {
	var mntDevices []types.MntDevice

	name := "Create Files on the Root Filesystem"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices = nil
	config := `{
	  "ignition": { "version": "2.0.0" },
	  "storage": {
	    "files": [{
	      "filesystem": "root",
	      "path": "/foo/bar",
	      "contents": { "source": "data:,example%20file%0A" }
	    }]
	  }
	}`
	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Name: "bar",
				Path: "foo",
			},
			Contents: []string{"example file\n"},
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}

func User_group_by_id_2_0_0() types.Test {
	var mntDevices []types.MntDevice

	name := "2.0.0 User/Group by id"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices = nil
	config := `{
	  "ignition": { "version": "2.0.0" },
	  "storage": {
	    "files": [{
	      "filesystem": "root",
	      "path": "/foo/bar",
	      "contents": { "source": "data:,example%20file%0A" },
		  "user": {"id": 500},
		  "group": {"id": 500}
	    }]
	  }
	}`
	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Name:  "bar",
				Path:  "foo",
				User:  500,
				Group: 500,
			},
			Contents: []string{"example file\n"},
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}

func User_group_by_id_2_1_0() types.Test {
	var mntDevices []types.MntDevice

	name := "2.1.0 User/Group by id"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices = nil
	config := `{
	  "ignition": { "version": "2.0.0" },
	  "storage": {
	    "files": [{
	      "filesystem": "root",
	      "path": "/foo/bar",
	      "contents": { "source": "data:,example%20file%0A" },
		  "user": {"id": 500},
		  "group": {"id": 500}
	    }]
	  }
	}`
	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Name:  "bar",
				Path:  "foo",
				User:  500,
				Group: 500,
			},
			Contents: []string{"example file\n"},
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}

func User_group_by_name_2_1_0() types.Test {
	var mntDevices []types.MntDevice

	name := "2.1.0 User/Group by name"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices = nil
	config := `{
	  "ignition": { "version": "2.0.0" },
	  "storage": {
	    "files": [{
	      "filesystem": "root",
	      "path": "/foo/bar",
	      "contents": { "source": "data:,example%20file%0A" },
		  "user": {"name": "core"},
		  "group": {"name": "core"}
	    }]
	  }
	}`
	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Name:  "bar",
				Path:  "foo",
				User:  500,
				Group: 500,
			},
			Contents: []string{"example file\n"},
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}
