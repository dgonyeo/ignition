// Copyright 2016 CoreOS, Inc.
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

package validate

import (
	"errors"
	"reflect"
	"testing"

	"github.com/coreos/go-semver/semver"

	// Import into the same namespace to keep config definitions clean
	. "github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/config/validate/report"
)

func TestValidate(t *testing.T) {
	type in struct {
		cfg Config
	}
	type out struct {
		err error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{cfg: Config{Ignition: IgnitionVersion{Version: semver.Version{Major: 2}.String()}}},
			out: out{},
		},
		{
			in:  in{cfg: Config{}},
			out: out{err: ErrOldVersion},
		},
		{
			in: in{cfg: Config{
				Ignition: IgnitionVersion{
					Version: semver.Version{Major: 2}.String(),
					Config: IgnitionConfig{
						Replace: &ConfigReference{
							Verification: Verification{
								Hash: func(s string) *string { return &s }("foobar-"),
							},
						},
					},
				},
			}},
			out: out{errors.New("unrecognized hash function")},
		},
		{
			in: in{cfg: Config{
				Ignition: IgnitionVersion{Version: semver.Version{Major: 2}.String()},
				Storage: Storage{
					Filesystems: []Filesystem{
						{
							Name: "filesystem1",
							Mount: &Mount{
								Device: "/dev/disk/by-partlabel/ROOT",
								Format: "btrfs",
							},
						},
					},
				},
			}},
			out: out{},
		},
		{
			in: in{cfg: Config{
				Ignition: IgnitionVersion{Version: semver.Version{Major: 2}.String()},
				Storage: Storage{
					Filesystems: []Filesystem{
						{
							Name: "filesystem1",
							Path: func(p string) *string { return &p }("/sysroot"),
						},
					},
				},
			}},
			out: out{},
		},
		{
			in: in{cfg: Config{
				Ignition: IgnitionVersion{Version: semver.Version{Major: 2}.String()},
				Systemd:  Systemd{Units: []Unit{{Name: "foo.bar", Contents: "[Foo]\nfoo=qux"}}},
			}},
			out: out{err: errors.New("invalid systemd unit extension")},
		},
		{
			in: in{cfg: Config{
				Ignition: IgnitionVersion{Version: semver.Version{Major: 2}.String()},
				Networkd: Networkd{Units: []NetworkdUnit{{Name: "foo.link", Contents: ""}}},
			}},
			out: out{err: errors.New("invalid or empty unit content")},
		},
	}

	for i, test := range tests {
		r := ValidateWithoutSource(reflect.ValueOf(test.in.cfg))
		expectedReport := report.ReportFromError(test.out.err, report.EntryError)
		if !reflect.DeepEqual(expectedReport, r) {
			t.Errorf("#%d: bad error: want %v, got %v", i, expectedReport, r)
		}
	}
}
