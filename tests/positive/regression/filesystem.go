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

package regression

import (
	"github.com/coreos/ignition/tests/types"
)

func Equivalent_filesystem_uuids_treated_distinct_ext4() types.Test {
	name := "Regression: Equivalent Filesystem UUIDs treated as distinct - ext4"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices := []types.MntDevice{
		{
			Label: "EFI-SYSTEM",
			Code:  "$DEVICE",
		},
	}
	config := `{
		"ignition": {"version": "2.1.0"},
		"storage": {
		    "filesystems": [
		      {
		        "mount": {
		          "device": "$DEVICE",
		          "format": "ext4",
		          "uuid": "6ABE925E-6DAF-4FAD-BC09-8D56BE8822DE"
		        }
		      }
		    ]
		  }
		}`
	in[0].Partitions.GetPartition("EFI-SYSTEM").FilesystemUUID = "6ABE925E-6DAF-4FAD-BC09-8D56BE8822DE"
	out[0].Partitions.GetPartition("EFI-SYSTEM").FilesystemUUID = "6ABE925E-6DAF-4FAD-BC09-8D56BE8822DE"

	return types.Test{name, in, out, mntDevices, config}
}

func Equivalent_filesystem_uuids_treated_distinct_vfat() types.Test {
	name := "Regression: Equivalent Filesystem UUIDs treated as distinct - vfat"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices := []types.MntDevice{
		{
			Label: "EFI-SYSTEM",
			Code:  "$DEVICE",
		},
	}
	config := `{
		"ignition": {"version": "2.1.0"},
		"storage": {
		    "filesystems": [
		      {
		        "mount": {
		          "device": "$DEVICE",
		          "format": "vfat",
		          "uuid": "2E24EC82"
		        }
		      }
		    ]
		  }
		}`
	in[0].Partitions.GetPartition("EFI-SYSTEM").FilesystemUUID = "2e24ec82"
	in[0].Partitions.GetPartition("EFI-SYSTEM").FilesystemType = "vfat"
	out[0].Partitions.GetPartition("EFI-SYSTEM").FilesystemUUID = "2e24ec82"
	out[0].Partitions.GetPartition("EFI-SYSTEM").FilesystemType = "vfat"

	return types.Test{name, in, out, mntDevices, config}
}
