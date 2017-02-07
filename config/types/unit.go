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
	"fmt"
	"path/filepath"

	"github.com/coreos/go-systemd/unit"

	"github.com/coreos/ignition/config/validate/report"
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
			Message: fmt.Sprintf("invalid systemd unit extension: %q", filepath.Ext(u.Name)),
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
			Message: fmt.Sprintf("invalid networkd unit extension: %q", filepath.Ext(u.Name)),
			Kind:    report.EntryError,
		})
	}

	return r
}

//type SystemdUnit struct {
//	Name     SystemdUnitName     `json:"name,omitempty"`
//	Enable   bool                `json:"enable,omitempty"`
//	Mask     bool                `json:"mask,omitempty"`
//	Contents string              `json:"contents,omitempty"`
//	DropIns  []SystemdUnitDropIn `json:"dropins,omitempty"`
//}

//func (u SystemdUnit) Validate() report.Report {
//	if err := validateUnitContent(u.Contents); err != nil {
//		if err != errEmptyUnit || (err == errEmptyUnit && len(u.DropIns) == 0) {
//			return report.ReportFromError(err, report.EntryError)
//		}
//	}
//
//	return report.Report{}
//}

//type SystemdUnitDropIn struct {
//	Name     SystemdUnitDropInName `json:"name,omitempty"`
//	Contents string                `json:"contents,omitempty"`
//}
//
//func (u SystemdUnitDropIn) Validate() report.Report {
//	if err := validateUnitContent(u.Contents); err != nil {
//		return report.ReportFromError(err, report.EntryError)
//	}
//
//	return report.Report{}
//}
//
//type SystemdUnitName string
//
//func (n SystemdUnitName) Validate() report.Report {
//	switch filepath.Ext(string(n)) {
//	case ".service", ".socket", ".device", ".mount", ".automount", ".swap", ".target", ".path", ".timer", ".snapshot", ".slice", ".scope":
//		return report.Report{}
//	default:
//		return report.ReportFromError(errors.New("invalid systemd unit extension"), report.EntryError)
//	}
//}
//
//type SystemdUnitDropInName string
//
//func (n SystemdUnitDropInName) Validate() report.Report {
//	switch filepath.Ext(string(n)) {
//	case ".conf":
//		return report.Report{}
//	default:
//		return report.ReportFromError(errors.New("invalid systemd unit drop-in extension"), report.EntryError)
//	}
//}

//type NetworkdUnit struct {
//	Name     NetworkdUnitName `json:"name,omitempty"`
//	Contents string           `json:"contents,omitempty"`
//}
//
//func (u NetworkdUnit) Validate() report.Report {
//	if err := validateUnitContent(u.Contents); err != nil {
//		return report.ReportFromError(err, report.EntryError)
//	}
//
//	return report.Report{}
//}
//
//type NetworkdUnitName string
//
//func (n NetworkdUnitName) Validate() report.Report {
//	switch filepath.Ext(string(n)) {
//	case ".link", ".netdev", ".network":
//		return report.Report{}
//	default:
//		return report.ReportFromError(errors.New("invalid networkd unit extension"), report.EntryError)
//	}
//}

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
