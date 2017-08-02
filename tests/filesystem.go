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
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/coreos/ignition/tests/types"
)

func prepareRootPartitionForPasswd(t *testing.T, partitions []*types.Partition) {
	for _, p := range partitions {
		if p.Label == "ROOT" {
			_ = os.MkdirAll(filepath.Join(p.MountPath, "home"), 0755)
			_ = os.MkdirAll(filepath.Join(p.MountPath, "usr", "bin"), 0755)
			_ = os.MkdirAll(filepath.Join(p.MountPath, "usr", "sbin"), 0755)
			_ = os.MkdirAll(filepath.Join(p.MountPath, "usr", "lib64"), 0755)
			_ = os.MkdirAll(filepath.Join(p.MountPath, "etc"), 0755)

			_ = os.Symlink(
				filepath.Join(p.MountPath, "usr", "lib64"), filepath.Join(
					p.MountPath, "lib64"))
			_ = os.Symlink(
				filepath.Join(p.MountPath, "usr", "bin"), filepath.Join(
					p.MountPath, "bin"))
			_ = os.Symlink(
				filepath.Join(p.MountPath, "usr", "sbin"), filepath.Join(
					p.MountPath, "sbin"))

			_, _ = exec.Command(
				"cp", "/etc/ld.so.cache", filepath.Join(
					p.MountPath, "etc")).CombinedOutput()
			_, _ = exec.Command(
				"cp", "/lib64/libblkid.so.1", filepath.Join(
					p.MountPath, "usr", "lib64")).CombinedOutput()
			_, _ = exec.Command(
				"cp", "/lib64/libpthread.so.0", filepath.Join(
					p.MountPath, "usr", "lib64")).CombinedOutput()
			_, _ = exec.Command(
				"cp", "/lib64/libc.so.6", filepath.Join(
					p.MountPath, "usr", "lib64")).CombinedOutput()
			_, _ = exec.Command(
				"cp", "/lib64/libuuid.so.1", filepath.Join(
					p.MountPath, "usr", "lib64")).CombinedOutput()
			_, _ = exec.Command(
				"cp", "/lib64/ld-linux-x86-64.so.2", filepath.Join(
					p.MountPath, "usr", "lib64")).CombinedOutput()
			_, _ = exec.Command(
				"cp", "/lib64/libnss_files.so.2", filepath.Join(
					p.MountPath, "usr", "lib64")).CombinedOutput()

			_, _ = exec.Command(
				"cp", "bin/amd64/id", filepath.Join(
					p.MountPath, "usr", "bin")).CombinedOutput()
			_, _ = exec.Command(
				"cp", "bin/amd64/useradd", filepath.Join(
					p.MountPath, "usr", "sbin")).CombinedOutput()
			_, _ = exec.Command(
				"cp", "bin/amd64/usermod", filepath.Join(
					p.MountPath, "usr", "sbin")).CombinedOutput()
		}
	}
}

func getRootLocation(partitions []*types.Partition) string {
	for _, p := range partitions {
		if p.Label == "ROOT" {
			return p.MountPath
		}
	}
	return ""
}

func removeFile(t *testing.T, imageFile string) {
	err := os.Remove(imageFile)
	if err != nil {
		t.Log(err)
	}
}

func removeMountFolders(t *testing.T, partitions []*types.Partition) {
	for _, p := range partitions {
		err := os.RemoveAll(p.MountPath)
		if err != nil {
			t.Log(err)
		}
	}
}

// returns true if no error, false if error
func runIgnition(t *testing.T, stage, root, configDir string, expectFail bool) bool {
	cmd := exec.Command(
		"ignition", "-clear-cache", "-oem",
		"file", "-stage", stage, "-root", root)
	t.Log("ignition", "-clear-cache", "-oem", "file", "-stage", stage, "-root", root)
	cmd.Dir = configDir
	out, err := cmd.CombinedOutput()
	if err != nil && !expectFail {
		t.Fatal("ignition", err, string(out))
	}
	return err == nil
}

// pickDevice will return the device corresponding to a partition with a given
// label in the given imageFile
func pickDevice(t *testing.T, partitions []*types.Partition, fileName string, label string) string {
	number := -1
	for _, p := range partitions {
		if p.Label == label {
			number = p.Number
		}
	}
	if number == -1 {
		return ""
	}

	kpartxOut, err := exec.Command("kpartx", "-l", fileName).CombinedOutput()
	if err != nil {
		t.Fatal("kpartx -l", err, string(kpartxOut))
	}
	re, err := regexp.Compile("/dev/(?P<device>[\\w\\d]+)")
	if err != nil {
		t.Fatal("compiling regexp:", err)
	}
	match := re.FindSubmatch(kpartxOut)
	if len(match) < 2 {
		t.Log(string(kpartxOut))
		t.Fatal("couldn't find device")
	}
	return fmt.Sprintf("/dev/mapper/%sp%d", string(match[1]), number)
}

func writeIgnitionConfig(t *testing.T, config string) string {
	tmpDir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(
		filepath.Join(tmpDir, "config.ign"), []byte(config), 0644)
	if err != nil {
		t.Fatal(err)
	}
	return tmpDir
}

func calculateImageSize(partitions []*types.Partition) int64 {
	// 63 is the number of sectors cpgt uses when generating a hybrid MBR
	size := int64(63 * 512)
	for _, p := range partitions {
		size += int64(align(p.Length, 512) * 512)
	}
	size = size + int64(4096*512) // extra room to allow for alignments
	return size
}

// createVolume will create the image file of the specified size, create a
// partition table in it, and generate mount paths for every partition
func createVolume(t *testing.T, imageFile string, size int64, cylinders int, heads int, sectorsPerTrack int, partitions []*types.Partition) {
	// attempt to create the file, will leave already existing files alone.
	// os.Truncate requires the file to already exist
	out, err := os.Create(imageFile)
	if err != nil {
		t.Fatal("create", err, out)
	}
	out.Close()

	// Truncate the file to the given size
	err = os.Truncate(imageFile, size)
	if err != nil {
		t.Fatal("truncate", err)
	}

	createPartitionTable(t, imageFile, partitions)

	for counter, partition := range partitions {
		if partition.TypeCode == "blank" || partition.FilesystemType == "" {
			continue
		}

		mntPath, err := ioutil.TempDir("", fmt.Sprintf("hd1p%d", counter))
		if err != nil {
			t.Fatal(err)
		}
		partition.MountPath = mntPath
	}
}

// setDevices will create devices for each of the partitions in the imageFile,
// and then will format each partition according to what's descrived in the
// partitions argument.
func setDevices(t *testing.T, imageFile string, partitions []*types.Partition) string {
	loopDevice := kpartxAdd(t, imageFile)

	for _, partition := range partitions {
		if partition.TypeCode == "blank" || partition.FilesystemType == "" {
			continue
		}

		partition.Device = fmt.Sprintf(
			"/dev/mapper/%sp%d", loopDevice, partition.Number)
		formatPartition(t, partition)
	}
	return fmt.Sprintf("/dev/%s", loopDevice)
}

func destroyDevices(t *testing.T, imageFile string) {
	kpartxOut, err := exec.Command(
		"kpartx", "-d", imageFile).CombinedOutput()
	if err != nil {
		t.Fatal("kpartx", err, string(kpartxOut))
	}
}
func formatPartition(t *testing.T, partition *types.Partition) {
	switch partition.FilesystemType {
	case "vfat":
		formatVFAT(t, partition)
	case "ext2", "ext4":
		formatEXT(t, partition)
	case "btrfs":
		formatBTRFS(t, partition)
	case "xfs":
		formatXFS(t, partition)
	case "swap":
		formatSWAP(t, partition)
	default:
		if partition.FilesystemType == "blank" ||
			partition.FilesystemType == "" {
			return
		}
		t.Fatal("Unknown partition", partition.FilesystemType)
	}
}

func formatSWAP(t *testing.T, partition *types.Partition) {
	opts := []string{}
	if partition.FilesystemLabel != "" {
		opts = append(opts, "-L", partition.FilesystemLabel)
	}
	if partition.FilesystemUUID != "" {
		opts = append(opts, "-U", partition.FilesystemUUID)
	}
	opts = append(
		opts, partition.Device)
	out, err := exec.Command("mkswap", opts...).CombinedOutput()
	if err != nil {
		t.Fatal("mkswap", err, string(out))
	}
}

func formatVFAT(t *testing.T, partition *types.Partition) {
	opts := []string{}
	if partition.FilesystemLabel != "" {
		opts = append(opts, "-n", partition.FilesystemLabel)
	}
	if partition.FilesystemUUID != "" {
		opts = append(opts, "-i", partition.FilesystemUUID)
	}
	opts = append(
		opts, partition.Device)
	out, err := exec.Command("mkfs.vfat", opts...).CombinedOutput()
	if err != nil {
		t.Fatal("mkfs.vfat", err, string(out))
	}
}

func formatXFS(t *testing.T, partition *types.Partition) {
	opts := []string{}
	if partition.FilesystemLabel != "" {
		opts = append(opts, "-L", partition.FilesystemLabel)
	}
	if partition.FilesystemUUID != "" {
		opts = append(opts, "-m", partition.FilesystemUUID)
	}
	opts = append(
		opts, partition.Device)
	out, err := exec.Command("mkfs.xfs", opts...).CombinedOutput()
	if err != nil {
		t.Fatal("mkfs.xfs", err, string(out))
	}
}

func formatEXT(t *testing.T, partition *types.Partition) {
	out, err := exec.Command(
		"mke2fs", "-q", "-t", partition.FilesystemType, "-b", "4096",
		"-i", "4096", "-I", "128", partition.Device).CombinedOutput()
	if err != nil {
		t.Fatal("mke2fs", err, string(out))
	}

	opts := []string{"-e", "remount-ro"}
	if partition.FilesystemLabel != "" {
		opts = append(opts, "-L", partition.FilesystemLabel)
	}
	if partition.FilesystemUUID != "" {
		opts = append(opts, "-U", partition.FilesystemUUID)
	}

	if partition.TypeCode == "coreos-usr" {
		opts = append(
			opts, "-U", "clear", "-T", "20091119110000", "-c", "0", "-i", "0",
			"-m", "0", "-r", "0")
	}
	opts = append(opts, partition.Device)
	tuneOut, err := exec.Command("tune2fs", opts...).CombinedOutput()
	if err != nil {
		t.Fatal("tune2fs", err, string(tuneOut))
	}
}

func formatBTRFS(t *testing.T, partition *types.Partition) {
	opts := []string{}
	if partition.FilesystemLabel != "" {
		opts = append(opts, "--label", partition.FilesystemLabel)
	}
	if partition.FilesystemUUID != "" {
		opts = append(opts, "--uuid", partition.FilesystemUUID)
	}
	opts = append(opts, partition.Device)
	out, err := exec.Command("mkfs.btrfs", opts...).CombinedOutput()
	if err != nil {
		t.Fatal("mkfs.btrfs", err, string(out))
	}
}

func align(count int, alignment int) int {
	offset := count % alignment
	if offset != 0 {
		count += alignment - offset
	}
	return count
}

func setOffsets(partitions []*types.Partition) {
	// 34 is the first non-reserved GPT sector
	offset := 34
	for _, p := range partitions {
		if p.Length == 0 || p.TypeCode == "blank" {
			continue
		}
		offset = align(offset, 4096)
		p.Offset = offset
		offset += p.Length
	}
}

func createPartitionTable(t *testing.T, imageFile string, partitions []*types.Partition) {
	opts := []string{imageFile}
	hybrids := []int{}
	for _, p := range partitions {
		if p.TypeCode == "blank" || p.Length == 0 {
			continue
		}
		opts = append(opts, fmt.Sprintf(
			"--new=%d:%d:+%d", p.Number, p.Offset, p.Length))
		opts = append(opts, fmt.Sprintf(
			"--change-name=%d:%s", p.Number, p.Label))
		if p.TypeGUID != "" {
			opts = append(opts, fmt.Sprintf(
				"--typecode=%d:%s", p.Number, p.TypeGUID))
		}
		if p.GUID != "" {
			opts = append(opts, fmt.Sprintf(
				"--partition-guid=%d:%s", p.Number, p.GUID))
		}
		if p.Hybrid {
			hybrids = append(hybrids, p.Number)
		}
	}
	if len(hybrids) > 0 {
		if len(hybrids) > 3 {
			t.Fatal("Can't have more than three hybrids")
		} else {
			opts = append(opts, fmt.Sprintf("-h=%s", intJoin(hybrids, ":")))
		}
	}
	sgdiskOut, err := exec.Command("sgdisk", opts...).CombinedOutput()
	if err != nil {
		t.Fatal("sgdisk", err, string(sgdiskOut))
	}
}

// kpartxAdd will use kpartx to add partition mappings for all of the partitions
// contained in the imageFile. This creates devices such as /dev/mapper/loop1p1
// corresponding to partitions in the imageFile. The loop device name (e.g.
// "loop1") will be returned.
func kpartxAdd(t *testing.T, imageFile string) string {
	kpartxOut, err := exec.Command(
		"kpartx", "-avs", imageFile).CombinedOutput()
	if err != nil {
		t.Fatal("kpartx", err, string(kpartxOut))
	}
	kpartxOut, err = exec.Command(
		"kpartx", "-l", imageFile).CombinedOutput()
	if err != nil {
		t.Fatal(err, string(kpartxOut))
	}
	re, err := regexp.Compile("/dev/(?P<device>[\\w\\d]+)")
	if err != nil {
		t.Fatal("compiling regexp:", err)
	}
	match := re.FindSubmatch(kpartxOut)
	if len(match) < 2 {
		t.Log(string(kpartxOut))
		t.Fatal("couldn't find device")
	}
	return string(match[1])
}

func mountRootPartition(t *testing.T, partitions []*types.Partition) bool {
	for _, partition := range partitions {
		if partition.Label != "ROOT" {
			continue
		}
		mountOut, err := exec.Command(
			"mount", partition.Device,
			partition.MountPath).CombinedOutput()
		if err != nil {
			t.Fatal("mount", err, string(mountOut))
		}
		return true
	}
	return false
}

func mountPartitions(t *testing.T, partitions []*types.Partition) {
	for _, partition := range partitions {
		if partition.FilesystemType == "" || partition.Label == "ROOT" {
			continue
		}
		mountOut, err := exec.Command(
			"mount", partition.Device,
			partition.MountPath).CombinedOutput()
		if err != nil {
			t.Fatal("mount", err, string(mountOut))
		}
	}
}

func updateTypeGUID(t *testing.T, partition *types.Partition) {
	switch partition.TypeCode {
	case "coreos-resize":
		partition.TypeGUID = "3884DD41-8582-4404-B9A8-E9B84F2DF50E"
	case "data":
		partition.TypeGUID = "0FC63DAF-8483-4772-8E79-3D69D8477DE4"
	case "coreos-rootfs":
		partition.TypeGUID = "5DFBF5F4-2848-4BAC-AA5E-0D9A20B745A6"
	case "bios":
		partition.TypeGUID = "21686148-6449-6E6F-744E-656564454649"
	case "efi":
		partition.TypeGUID = "C12A7328-F81F-11D2-BA4B-00A0C93EC93B"
	case "coreos-reserved":
		partition.TypeGUID = "C95DC21A-DF0E-4340-8D7B-26CBFA9A03E0"
	case "", "blank":
		return
	default:
		t.Fatal("Unknown TypeCode", partition.TypeCode)
	}
}

func intJoin(ints []int, delimiter string) string {
	strArr := []string{}
	for _, i := range ints {
		strArr = append(strArr, strconv.Itoa(i))
	}
	return strings.Join(strArr, delimiter)
}

func removeEmpty(strings []string) []string {
	var r []string
	for _, str := range strings {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func generateUUID(t *testing.T) string {
	out, err := exec.Command("uuidgen").CombinedOutput()
	if err != nil {
		t.Fatal("uuidgen", err)
	}
	return strings.TrimSpace(string(out))
}

func createFiles(t *testing.T, partitions []*types.Partition) {
	for _, partition := range partitions {
		if partition.Files == nil {
			continue
		}
		for _, file := range partition.Files {
			err := os.MkdirAll(filepath.Join(
				partition.MountPath, file.Path), 0644)
			if err != nil {
				t.Fatal("mkdirall", err)
			}
			f, err := os.Create(filepath.Join(
				partition.MountPath, file.Path, file.Name))
			defer f.Close()
			if err != nil {
				t.Fatal("create", err, f)
			}
			if file.Contents != nil {
				writer := bufio.NewWriter(f)
				writeStringOut, err := writer.WriteString(
					strings.Join(file.Contents, "\n"))
				if err != nil {
					t.Fatal("writeString", err, string(writeStringOut))
				}
				writer.Flush()
			}
		}
	}
}

func unmountRootPartition(t *testing.T, partitions []*types.Partition) {
	for _, partition := range partitions {
		if partition.Label != "ROOT" {
			continue
		}

		umountOut, err := exec.Command(
			"umount", partition.Device).CombinedOutput()
		if err != nil {
			t.Fatal("umount", err, string(umountOut))
		}
	}
}

func unmountPartitions(t *testing.T, partitions []*types.Partition) {
	for _, partition := range partitions {
		if partition.FilesystemType == "" || partition.Label == "ROOT" {
			continue
		}
		umountOut, err := exec.Command(
			"umount", partition.Device).CombinedOutput()
		if err != nil {
			t.Fatal("umount", err, string(umountOut))
		}
	}
}

func setExpectedPartitionsDrive(actual []*types.Partition, expected []*types.Partition) {
	for _, a := range actual {
		for _, e := range expected {
			if a.Number == e.Number {
				e.MountPath = a.MountPath
				e.Device = a.Device
				break
			}
		}
	}
}
