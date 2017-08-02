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

package blackbox

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/coreos/ignition/internal/exec/util"
	"github.com/coreos/ignition/tests/types"
)

func validatePartitions(t *testing.T, expected []*types.Partition, imageFile string) {
	for _, e := range expected {
		if e.TypeCode == "blank" {
			continue
		}
		sgdiskInfo, err := exec.Command(
			"sgdisk", "-i", strconv.Itoa(e.Number),
			imageFile).CombinedOutput()
		if err != nil {
			t.Error("sgdisk -i", strconv.Itoa(e.Number), err)
			return
		}

		re, err := regexp.Compile("Partition unique GUID: (?P<partition_guid>[\\d\\w-]+)")
		if err != nil {
			t.Fatal("compiling regexp:", err)
		}
		match := re.FindSubmatch(sgdiskInfo)
		if len(match) < 2 {
			t.Fatal("couldn't find GUID")
		}
		actualGUID := string(match[1])

		re, err = regexp.Compile("Partition GUID code: (?P<partition_code>[\\d\\w-]+)")
		if err != nil {
			t.Fatal("compiling regexp:", err)
		}
		match = re.FindSubmatch(sgdiskInfo)
		if len(match) < 2 {
			t.Fatal("couldn't find type GUID")
		}
		actualTypeGUID := string(match[1])

		re, err = regexp.Compile("Partition size: (?P<sectors>\\d+) sectors")
		if err != nil {
			t.Fatal("compiling regexp:", err)
		}
		match = re.FindSubmatch(sgdiskInfo)
		if len(match) < 2 {
			t.Fatal("couldn't find partition size")
		}
		actualSectors := string(match[1])

		re, err = regexp.Compile("Partition name: '(?P<name>[\\d\\w-_]+)'")
		if err != nil {
			t.Fatal("compiling regexp:", err)
		}
		match = re.FindSubmatch(sgdiskInfo)
		if len(match) < 2 {
			t.Fatal("couldn't find partition name")
		}
		actualLabel := string(match[1])

		// have to align the size to the nearest sector first
		expectedSectors := align(e.Length, 512)

		if e.TypeGUID != "" && e.TypeGUID != actualTypeGUID {
			t.Error("TypeGUID does not match!", e.TypeGUID, actualTypeGUID)
		}
		if e.GUID != "" && formatUUID(e.GUID) != formatUUID(actualGUID) {
			t.Error("GUID does not match!", e.GUID, actualGUID)
		}
		if e.Label != actualLabel {
			t.Error("Label does not match!", e.Label, actualLabel)
		}
		if strconv.Itoa(expectedSectors) != actualSectors {
			t.Error(
				"Sectors does not match!", expectedSectors, actualSectors)
		}
	}
}

func formatUUID(s string) string {
	return strings.ToUpper(strings.Replace(s, "-", "", -1))
}

func validateFilesystems(t *testing.T, expected []*types.Partition, imageFile string) {
	for _, e := range expected {
		if e.FilesystemType != "" {
			filesystemType, err := util.FilesystemType(e.Device)
			if err != nil {
				t.Fatalf("couldn't determine filesystem type: %v", err)
			}
			if filesystemType != e.FilesystemType {
				t.Errorf("FilesystemType does not match, expected:%q actual:%q",
					e.FilesystemType, filesystemType)
			}
		}
		if e.FilesystemUUID != "" {
			filesystemUUID, err := util.FilesystemUUID(e.Device)
			if err != nil {
				t.Fatalf("couldn't determine filesystem uuid: %v", err)
			}
			if formatUUID(filesystemUUID) != formatUUID(e.FilesystemUUID) {
				t.Errorf("FilesystemUUID does not match, expected:%q actual:%q",
					e.FilesystemUUID, filesystemUUID)
			}
		}
		if e.FilesystemLabel != "" {
			filesystemLabel, err := util.FilesystemLabel(e.Device)
			if err != nil {
				t.Fatalf("couldn't determine filesystem label: %v", err)
			}
			if filesystemLabel != e.FilesystemLabel {
				t.Errorf("FilesystemLabel does not match, expected:%q actual:%q",
					e.FilesystemLabel, filesystemLabel)
			}
		}
	}
}

func validateFilesDirectoriesAndLinks(t *testing.T, expected []*types.Partition) {
	for _, partition := range expected {
		for _, file := range partition.Files {
			validateFile(t, partition, file)
		}
		for _, dir := range partition.Directories {
			validateDirectory(t, partition, dir)
		}
		for _, link := range partition.Links {
			validateLink(t, partition, link)
		}
		for _, node := range partition.RemovedNodes {
			path := filepath.Join(partition.MountPath, node.Path, node.Name)
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Error("Node was expected to be removed and is present!", path)
			}
		}
	}
}

func validateFile(t *testing.T, partition *types.Partition, file types.File) {
	path := filepath.Join(partition.MountPath, file.Node.Path, file.Node.Name)
	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Errorf("Error stat'ing file %s: %v", path, err)
		return
	}
	if file.Contents != nil {
		expectedContents := strings.Join(file.Contents, "\n")
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			t.Error("Error when reading file", path)
			return
		}

		actualContents := string(dat)
		if expectedContents != actualContents {
			t.Error("Contents of file", path, "do not match!",
				expectedContents, actualContents)
		}
	}

	validateNode(t, fileInfo, file.Node)
}

func validateDirectory(t *testing.T, partition *types.Partition, dir types.Directory) {
	path := filepath.Join(partition.MountPath, dir.Node.Path, dir.Node.Name)
	dirInfo, err := os.Stat(path)
	if err != nil {
		t.Errorf("Error stat'ing directory %s: %v", path, err)
		return
	}
	if !dirInfo.IsDir() {
		t.Errorf("Node at %s is not a directory!", path)
	}
	validateMode(t, path, dir.Mode)
	validateNode(t, dirInfo, dir.Node)
}

func validateLink(t *testing.T, partition *types.Partition, link types.Link) {
	linkPath := filepath.Join(partition.MountPath, link.Node.Path, link.Node.Name)
	linkInfo, err := os.Lstat(linkPath)
	if err != nil {
		t.Error("Error stat'ing link \"" + linkPath + "\": " + err.Error())
		return
	}
	if link.Hard {
		targetPath := filepath.Join(partition.MountPath, link.Target)
		targetInfo, err := os.Stat(targetPath)
		if err != nil {
			t.Error("Error stat'ing target \"" + targetPath + "\": " + err.Error())
			return
		}
		if linkInfoSys, ok := linkInfo.Sys().(*syscall.Stat_t); ok {
			if targetInfoSys, ok := targetInfo.Sys().(*syscall.Stat_t); ok {
				if linkInfoSys.Ino != targetInfoSys.Ino {
					t.Error("Hard link and target don't have same inode value: " + linkPath + " " + targetPath)
					return
				}
			} else {
				t.Error("Stat type assertion failed, this will only work on Linux")
				return
			}
		} else {
			t.Error("Stat type assertion failed, this will only work on Linux")
			return
		}
	} else {
		if linkInfo.Mode()&os.ModeSymlink == 0 {
			t.Errorf("Node at symlink path is not a symlink (it's a %s): %s", linkInfo.Mode().String(), linkPath)
			return
		}
		targetPath, err := os.Readlink(linkPath)
		if err != nil {
			t.Error("Error reading symbolic link: " + err.Error())
			return
		}
		if targetPath != link.Target {
			t.Error("Actual and expected symbolic link targets don't match")
			return
		}
	}
	validateNode(t, linkInfo, link.Node)
}

func validateMode(t *testing.T, path string, mode string) {
	if mode != "" {
		sout, err := exec.Command(
			"stat", "-c", "%a", path).CombinedOutput()
		if err != nil {
			t.Error("Error running stat on node", path, err)
			return
		}
		statOut := strings.TrimSpace(string(sout))
		if mode != statOut {
			t.Error(
				"Node Mode does not match", path, mode, statOut)
		}
	}
}

func validateNode(t *testing.T, nodeInfo os.FileInfo, node types.Node) {
	if nodeInfoSys, ok := nodeInfo.Sys().(*syscall.Stat_t); ok {
		if nodeInfoSys.Uid != uint32(node.User) {
			t.Error("Node has the wrong owner", node.User, nodeInfoSys.Uid)
		}

		if nodeInfoSys.Gid != uint32(node.Group) {
			t.Error("Node has the wrong group owner", node.Group, nodeInfoSys.Gid)
		}
	} else {
		t.Error("Stat type assertion failed, this will only work on Linux")
		return
	}
}
