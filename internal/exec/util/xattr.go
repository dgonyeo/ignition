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

package util

import (
	"golang.org/x/sys/unix"
)

func retrieveXattrData(f func(buf []byte) (int, error)) ([]byte, error) {
	size, err := f(nil)
	if err != nil {
		return nil, err
	}
	if size > 0 {
		data := make([]byte, size)
		read, err := f(data)
		if err != nil {
			return nil, err
		}
		return data[:read], nil
	}
	return []byte{}, nil
}

func ListXattrs(path string) ([]string, error) {
	buf, err := retrieveXattrData(func(buf []byte) (int, error) {
		return unix.Listxattr(path, buf)
	})
	if err != nil {
		return nil, err
	}
	var result []string
	offset := 0
	for index, b := range buf {
		if b == 0 {
			result = append(result, string(buf[offset:index]))
			offset = index + 1
		}
	}
	return result, nil
}

func GetXattr(path, name string) ([]byte, error) {
	return retrieveXattrData(func(buf []byte) (int, error) {
		return unix.Getxattr(path, name, buf)
	})
}

func SetXattr(path, name string, data []byte) error {
	return unix.Setxattr(path, name, data, 0)
}

func LsetXattr(path, name string, data []byte) error {
	return unix.Lsetxattr(path, name, data, 0)
}
