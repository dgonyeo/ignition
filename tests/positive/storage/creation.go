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

package storage

import (
	"github.com/coreos/ignition/tests/types"
)

func Force_new_filesystem_of_same_type() types.Test {
	name := "Force new Filesystem Creation of same type"
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
					}}
				 }]
			}
	}`

	out[0].Partitions.GetPartition("EFI-SYSTEM").Files = []types.File{}
	out[0].Partitions.AddRemovedNodes("EFI-SYSTEM", []types.Node{
		{
			Name: "multiLine",
			Path: "path/example",
		}, {
			Name: "singleLine",
			Path: "another/path/example",
		}, {
			Name: "emptyFile",
			Path: "empty",
		}, {
			Name: "noPath",
			Path: "",
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}

func Wipe_filesystem_with_same_type() types.Test {
	name := "Wipe Filesystem with Filesystem of same type"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices := []types.MntDevice{
		{
			Label: "EFI-SYSTEM",
			Code:  "$DEVICE",
		},
	}
	config := `{
		"ignition": { "version": "2.1.0" },
		"storage": {
			"filesystems": [{
				"mount": {
					"device": "$DEVICE",
					"format": "ext4",
					"wipeFilesystem": true
				}}]
			}
	}`

	out[0].Partitions.GetPartition("EFI-SYSTEM").Files = []types.File{}
	out[0].Partitions.AddRemovedNodes("EFI-SYSTEM", []types.Node{
		{
			Name: "multiLine",
			Path: "path/example",
		}, {
			Name: "singleLine",
			Path: "another/path/example",
		}, {
			Name: "emptyFile",
			Path: "empty",
		}, {
			Name: "noPath",
			Path: "",
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}

func Create_new_partitions() types.Test {
	var mntDevices []types.MntDevice

	name := "Create new partitions"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices = nil
	config := `{
		"ignition": {"version": "2.1.0"},
		"storage": {
		    "disks": [
			    {
					"device": "$blackbox_ignition_secondary_disk.img",
					"wipeTable": true,
					"partitions": [
						{
							"label": "important-data",
							"number": 1,
							"size": 65536,
							"typeGuid": "B921B045-1DF0-41C3-AF44-4C6F280D3FAE",
							"guid": "8A7A6E26-5E8F-4CCA-A654-46215D4696AC"
						},
						{
							"label": "ephemeral-data",
							"number": 2,
							"size": 131072,
							"typeGuid": "CA7D7CCB-63ED-4C53-861C-1742536059CC",
							"guid": "A51034E6-26B3-48DF-BEED-220562AC7AD1"
						}
					]
				}
			]
		}
	}`
	// We need dummy partitions to get the test to pass on Fedora (kpartx acts
	// weird during test setup without them), the UUIDs in the input partitions
	// are intentionally different so if Ignition doesn't do the right thing the
	// validation will fail.
	in = append(in, types.Disk{
		ImageFile: "blackbox_ignition_secondary_disk.img",
		Partitions: types.Partitions{
			{
				Label:    "important-data",
				Number:   1,
				Length:   65536,
				TypeGUID: "B921B045-1DF0-41C3-AF44-4C6F280D3FAE",
				GUID:     "B921B045-1DF0-41C3-AF44-4C6F280D3FAE",
			},
			{
				Label:    "ephemeral-data",
				Number:   2,
				Length:   131072,
				TypeGUID: "CA7D7CCB-63ED-4C53-861C-1742536059CC",
				GUID:     "B921B045-1DF0-41C3-AF44-4C6F280D3FAE",
			},
		},
	})
	out = append(out, types.Disk{
		ImageFile: "blackbox_ignition_secondary_disk.img",
		Partitions: types.Partitions{
			{
				Label:    "important-data",
				Number:   1,
				Length:   65536,
				TypeGUID: "B921B045-1DF0-41C3-AF44-4C6F280D3FAE",
				GUID:     "8A7A6E26-5E8F-4CCA-A654-46215D4696AC",
			},
			{
				Label:    "ephemeral-data",
				Number:   2,
				Length:   131072,
				TypeGUID: "CA7D7CCB-63ED-4C53-861C-1742536059CC",
				GUID:     "A51034E6-26B3-48DF-BEED-220562AC7AD1",
			},
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}
