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

package types

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/coreos/go-systemd/unit"

	"github.com/coreos/ignition/config/validate/report"
)

var (
	ErrInvalidSystemdExt  = errors.New("invalid systemd unit extension")
	ErrInvalidNetworkdExt = errors.New("invalid networkd unit extension")
)

func (u Unit) Validate() report.Report {
	r := report.Report{}

	if err := validateUnitContent(u.Contents); err != nil {
		if err != errEmptyUnit || (err == errEmptyUnit && len(u.Dropins) == 0) {
			r.Add(report.Entry{
				Message: err.Error(),
				Kind:    report.EntryError,
			})
		}
	}

	switch filepath.Ext(u.Name) {
	case ".service", ".socket", ".device", ".mount", ".automount", ".swap", ".target", ".path", ".timer", ".snapshot", ".slice", ".scope":
	default:
		r.Add(report.Entry{
			Message: ErrInvalidSystemdExt.Error(),
			Kind:    report.EntryError,
		})
	}

	return r
}

func (d Dropin) Validate() report.Report {
	r := report.Report{}

	if err := validateUnitContent(d.Contents); err != nil {
		r.Add(report.Entry{
			Message: err.Error(),
			Kind:    report.EntryError,
		})
	}

	switch filepath.Ext(d.Name) {
	case ".conf":
	default:
		r.Add(report.Entry{
			Message: fmt.Sprintf("invalid systemd unit drop-in extension: %q", filepath.Ext(d.Name)),
			Kind:    report.EntryError,
		})
	}

	return r
}

func (u Networkdunit) Validate() report.Report {
	r := report.Report{}

	if err := validateUnitContent(u.Contents); err != nil {
		r.Add(report.Entry{
			Message: err.Error(),
			Kind:    report.EntryError,
		})
	}

	switch filepath.Ext(u.Name) {
	case ".link", ".netdev", ".network":
	default:
		r.Add(report.Entry{
			Message: ErrInvalidNetworkdExt.Error(),
			Kind:    report.EntryError,
		})
	}

	return r
}

var errEmptyUnit = fmt.Errorf("invalid or empty unit content")

func validateUnitContent(content string) error {
	c := bytes.NewBufferString(content)
	unit, err := unit.Deserialize(c)
	if err != nil {
		return fmt.Errorf("invalid unit content: %s", err)
	}

	if len(unit) == 0 {
		return errEmptyUnit
	}

	return nil
}
