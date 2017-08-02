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

package negative_storage

import (
	"github.com/coreos/ignition/tests/types"
)

func No_device() types.Test {
	var out []types.Disk
	var mntDevices []types.MntDevice

	name := "No Device"
	in := types.GetBaseDisk()
	out = in
	mntDevices = nil
	config := `{
		"ignition": {"version": "2.1.0"},
		"storage": {
			"filesystems": [{
				"mount": {
					"format": "ext4"
				},
				"name": "foobar"
			}]
		}
	}`

	return types.Test{name, in, out, mntDevices, config}
}

func No_device_with_force() types.Test {
	var out []types.Disk
	var mntDevices []types.MntDevice

	name := "No Device w/ Force"
	in := types.GetBaseDisk()
	out = in
	mntDevices = nil
	config := `{
		"ignition": {"version": "2.0.0"},
		"storage": {
			"filesystems": [{
				"mount": {
					"format": "ext4",
					"create": {
						"force": true
					}
				},
				"name": "foobar",
			}]
		}
	}`

	return types.Test{name, in, out, mntDevices, config}
}

func No_device_with_wipe_filesystem_true() types.Test {
	var out []types.Disk
	var mntDevices []types.MntDevice

	name := "No Device w/ wipeFilesystem true"
	in := types.GetBaseDisk()
	out = in
	mntDevices = nil
	config := `{
		"ignition": {"version": "2.1.0"},
		"storage": {
			"filesystems": [{
				"mount": {
					"format": "ext4",
					"wipeFilesystem": true
				},
				"name": "foobar"
			}]
		}
	}`

	return types.Test{name, in, out, mntDevices, config}
}

func No_device_with_wipe_filesystem_false() types.Test {
	var out []types.Disk
	var mntDevices []types.MntDevice

	name := "No Device w/ wipeFilesystem false"
	in := types.GetBaseDisk()
	out = in
	mntDevices = nil
	config := `{
		"ignition": {"version": "2.1.0"},
		"storage": {
			"filesystems": [{
				"mount": {
					"format": "ext4",
					"wipeFilesystem": false
				},
				"name": "foobar"
			}]
		}
	}`

	return types.Test{name, in, out, mntDevices, config}
}
