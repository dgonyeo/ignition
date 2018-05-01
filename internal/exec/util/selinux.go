// Copyright 2018 CoreOS, Inc.
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

// +build linux

package util

// #cgo LDFLAGS: -lselinux
// #include <stdlib.h>
// #include <selinux/selinux.h>
import "C"

import (
	"fmt"
	"unsafe"
)

func SelinuxSetPolicyRoot(rootPath string) error {
	cRootPath := C.CString(rootPath)
	defer C.free(unsafe.Pointer(cRootPath))
	res := C.selinux_set_policy_root(cRootPath)
	if res != 0 {
		return fmt.Errorf("selinux_set_policy_root: failed to set policy root")
	}
	return nil
}
