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
	"errors"
	"fmt"
	"net/url"

	"github.com/vincent-petithory/dataurl"

	"github.com/coreos/ignition/config/validate/report"
)

var (
	ErrCompressionInvalid = errors.New("invalid compression method")
)

func (fc FileContents) Validate() report.Report {
	r := report.Report{}
	switch fc.Compression {
	case "", "gzip":
	default:
		r.Add(report.Entry{
			Message: ErrCompressionInvalid.Error(),
			Kind:    report.EntryError,
		})
	}
	u, err := url.Parse(fc.Source)
	if err != nil {
		r.Add(report.Entry{
			Message: fmt.Sprintf("invalid url: %q: ", fc.Source, err.Error()),
			Kind:    report.EntryError,
		})
	}

	switch u.Scheme {
	case "http", "https", "oem":

	case "data":
		if _, err := dataurl.DecodeString(fc.Source); err != nil {
			r.Add(report.Entry{
				Message: fmt.Sprintf("invalid data url: %v", err.Error()),
				Kind:    report.EntryError,
			})
		}
	default:
		r.Add(report.Entry{
			Message: fmt.Sprintf("invalid url scheme: %q", u.Scheme),
			Kind:    report.EntryError,
		})
	}
	return r
}
