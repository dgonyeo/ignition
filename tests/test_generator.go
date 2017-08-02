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
	"github.com/coreos/ignition/tests/negative/files"
	"github.com/coreos/ignition/tests/negative/general"
	"github.com/coreos/ignition/tests/negative/regression"
	"github.com/coreos/ignition/tests/negative/storage"
	ntimeouts "github.com/coreos/ignition/tests/negative/timeouts"
	"github.com/coreos/ignition/tests/positive/files"
	"github.com/coreos/ignition/tests/positive/general"
	"github.com/coreos/ignition/tests/positive/networkd"
	"github.com/coreos/ignition/tests/positive/passwd"
	"github.com/coreos/ignition/tests/positive/regression"
	"github.com/coreos/ignition/tests/positive/storage"
	"github.com/coreos/ignition/tests/positive/systemd"
	ptimeouts "github.com/coreos/ignition/tests/positive/timeouts"
	"github.com/coreos/ignition/tests/types"
)

func createNegativeTests() []types.Test {
	tests := []types.Test{}

	tests = append(tests, negative_files.Invalid_hash())
	tests = append(tests, negative_general.Replace_config_with_invalid_hash())
	tests = append(tests, negative_general.Append_config_with_invalid_hash())
	tests = append(tests, negative_general.Invalid_version())
	tests = append(tests, negative_regression.Vfat_ignores_wipe_filesystem())
	tests = append(tests, negative_storage.Invalid_filesystem())
	tests = append(tests, negative_storage.No_device())
	tests = append(tests, negative_storage.No_device_with_force())
	tests = append(tests, negative_storage.No_device_with_wipe_filesystem_true())
	tests = append(tests, negative_storage.No_device_with_wipe_filesystem_false())
	tests = append(tests, negative_storage.No_filesystem_type())
	tests = append(tests, negative_storage.No_filesystem_type_with_force())
	tests = append(tests, negative_storage.No_filesystem_type_with_wipe_filesystem())
	tests = append(tests, ntimeouts.Decrease_HTTP_Response_Headers_Timeout())

	return tests
}

func createTests() []types.Test {
	tests := []types.Test{}

	tests = append(tests, files.Create_directory_on_root())
	tests = append(tests, files.Create_file_on_root())
	tests = append(tests, files.User_group_by_id_2_0_0())
	tests = append(tests, files.User_group_by_id_2_1_0())
	// TODO: Investigate why ignition's C code hates our environment
	// tests = append(tests, files.User_group_by_name_2_1_0())
	tests = append(tests, files.Validate_file_hash_from_data_url())
	tests = append(tests, files.Validate_file_hash_from_http_url())
	tests = append(tests, files.Create_hard_link_on_root())
	tests = append(tests, files.Create_symlink_on_root())
	tests = append(tests, files.Create_file_from_remote_contents())
	tests = append(tests, general.Reformat_rootfs_and_write_file())
	tests = append(tests, general.Set_hostname())
	tests = append(tests, general.Replace_config_with_remote_config())
	tests = append(tests, general.Append_config_with_remote_config())
	tests = append(tests, general.Empty_userdata())
	tests = append(tests, networkd.Create_networkd_unit())
	tests = append(tests, passwd.Add_passwd_users())
	tests = append(tests, regression.Equivalent_filesystem_uuids_treated_distinct_ext4())
	tests = append(tests, regression.Equivalent_filesystem_uuids_treated_distinct_vfat())
	tests = append(tests, storage.Force_new_filesystem_of_same_type())
	tests = append(tests, storage.Wipe_filesystem_with_same_type())
	tests = append(tests, storage.Create_new_partitions())
	tests = append(tests, storage.Reuse_existing_filesystem())
	tests = append(tests, storage.Reformat_to_btrfs())
	tests = append(tests, storage.Reformat_to_xfs())
	tests = append(tests, systemd.Create_systemd_service())
	tests = append(tests, systemd.Modify_systemd_service())
	tests = append(tests, systemd.Mask_systemd_services())
	tests = append(tests, ptimeouts.Increase_HTTP_Response_Headers_Timeout())
	tests = append(tests, ptimeouts.Confirm_HTTP_Backoff_Works())

	return tests
}
