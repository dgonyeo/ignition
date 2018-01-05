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
	"fmt"

	"github.com/coreos/ignition/config/validate/report"
)

func (n Disk) Validate() report.Report {
	return report.Report{}
}

func (n Disk) ValidateDevice() report.Report {
	if len(n.Device) == 0 {
		return report.ReportFromError(fmt.Errorf("disk device is required"), report.EntryError)
	}
	if err := validatePath(string(n.Device)); err != nil {
		return report.ReportFromError(err, report.EntryError)
	}
	return report.Report{}
}

func (n Disk) ValidatePartitions() report.Report {
	r := report.Report{}
	if n.partitionNumbersCollide() {
		r.Add(report.Entry{
			Message: fmt.Sprintf("disk %q: partition numbers collide", n.Device),
			Kind:    report.EntryError,
		})
	}
	if n.partitionsOverlap() {
		r.Add(report.Entry{
			Message: fmt.Sprintf("disk %q: partitions overlap", n.Device),
			Kind:    report.EntryError,
		})
	}
	if n.partitionsMisaligned() {
		r.Add(report.Entry{
			Message: fmt.Sprintf("disk %q: partitions misaligned", n.Device),
			Kind:    report.EntryError,
		})
	}
	// Disks which have no errors at this point will likely succeed in sgdisk
	return r
}

// partitionNumbersCollide returns true if partition numbers in n.Partitions are not unique.
func (n Disk) partitionNumbersCollide() bool {
	m := map[int][]Partition{}
	for _, p := range n.Partitions {
		if p.Number != 0 {
			// a number of 0 means next available number, multiple devices can specify this
			m[p.Number] = append(m[p.Number], p)
		}
	}
	for _, n := range m {
		if len(n) > 1 {
			// TODO(vc): return information describing the collision for logging
			return true
		}
	}
	return false
}

// end returns the last sector of a partition.
func (p Partition) end() int {
	if p.GetSize() == 0 {
		// a size of 0 means "fill available", just return the start as the end for those.
		return p.GetStart()
	}
	return p.GetStart() + p.GetSize() - 1
}

// partitionsOverlap returns true if any explicitly dimensioned partitions overlap
func (n Disk) partitionsOverlap() bool {
	for _, p := range n.Partitions {
		// Starts of 0 are placed by sgdisk into the "largest available block" at that time.
		// We aren't going to check those for overlap since we don't have the disk geometry.
		if p.GetStart() == 0 {
			continue
		}

		for _, o := range n.Partitions {
			if p == o || o.GetStart() == 0 {
				continue
			}

			// is p.Start within o?
			if p.GetStart() >= o.GetStart() && p.GetStart() <= o.end() {
				return true
			}

			// is p.end() within o?
			if p.end() >= o.GetStart() && p.end() <= o.end() {
				return true
			}

			// do p.Start and p.end() straddle o?
			if p.GetStart() < o.GetStart() && p.end() > o.end() {
				return true
			}
		}
	}
	return false
}

// partitionsMisaligned returns true if any of the partitions don't start on a 2048-sector (1MiB) boundary.
func (n Disk) partitionsMisaligned() bool {
	for _, p := range n.Partitions {
		if (p.GetStart() & (2048 - 1)) != 0 {
			return true
		}
	}
	return false
}
