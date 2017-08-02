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

package general

import (
	"github.com/coreos/ignition/tests/types"
)

func Reformat_rootfs_and_write_file() types.Test {
	name := "Reformat rootfs to ext4 & drop file in /ignition/test"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices := []types.MntDevice{
		{
			Label: "EFI-SYSTEM",
			Code:  "$DEVICE",
		},
	}
	config := `{
		"ignition": {"version": "2.0.0"},
		"storage": {
			"filesystems": [{
				"mount": {
					"device": "$DEVICE",
					"format": "ext4",
					"create": {
						"force": true
					}},
				 "name": "test"}],
			"files": [{
				"filesystem": "test",
				"path": "/ignition/test",
				"contents": {"source": "data:,asdf"}
			}]}
	}`

	in[0].Partitions.GetPartition("EFI-SYSTEM").FilesystemType = "ext2"
	out[0].Partitions.GetPartition("EFI-SYSTEM").Files = []types.File{
		{
			Node: types.Node{
				Name: "test",
				Path: "ignition",
			},
			Contents: []string{"asdf"},
		},
	}

	return types.Test{name, in, out, mntDevices, config}
}

func Set_hostname() types.Test {
	var mntDevices []types.MntDevice

	name := "Setting the hostname"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices = nil
	config := `{
	  "ignition": { "version": "2.0.0" },
	  "storage": {
	    "files": [{
	      "filesystem": "root",
	      "path": "/etc/hostname",
	      "mode": 420,
	      "contents": { "source": "data:,core1" }
	    }]
	  }
	}`
	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Name: "hostname",
				Path: "etc",
			},
			Mode:     "644",
			Contents: []string{"core1"},
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}

func Replace_config_with_remote_config() types.Test {
	var mntDevices []types.MntDevice

	name := "Replacing the Config with a Remote Config"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices = nil
	config := `{
	  "ignition": {
	    "version": "2.0.0",
	    "config": {
	      "replace": {
	        "source": "http://127.0.0.1:8080/config",
			"verification": { "hash": "sha512-41d9a1593dd4cbcacc966dce574523ffe3780ec2710716fab28b46f0f24d20b5ec49f307a9e9d331af958e508f472f32135c740d1214c5f02fc36016b538e7ff" }
	      }
	    }
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

func Append_config_with_remote_config() types.Test {
	var mntDevices []types.MntDevice

	name := "Appending to the Config with a Remote Config"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices = nil
	config := `{
	  "ignition": {
	    "version": "2.0.0",
	    "config": {
	      "append": [{
	        "source": "http://127.0.0.1:8080/config",
			"verification": { "hash": "sha512-41d9a1593dd4cbcacc966dce574523ffe3780ec2710716fab28b46f0f24d20b5ec49f307a9e9d331af958e508f472f32135c740d1214c5f02fc36016b538e7ff" }
	      }]
	    }
	  },
      "storage": {
        "files": [{
          "filesystem": "root",
          "path": "/foo/bar2",
          "contents": { "source": "data:,another%20example%20file%0A" }
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
		{
			Node: types.Node{
				Name: "bar2",
				Path: "foo",
			},
			Contents: []string{"another example file\n"},
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}

func Empty_userdata() types.Test {
	var mntDevices []types.MntDevice

	name := "Empty Userdata"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices = nil
	config := `{
		"ignition": {"version": "2.1.0"}
	}`

	return types.Test{name, in, out, mntDevices, config}
}
