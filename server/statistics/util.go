// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package statistics

import (
	"github.com/montanaflynn/stats"
)

// RollingStats provides rolling statistics with specified window size.
// There are window size records for calculating.
type RollingStats struct {
	records []float64
	size    int
	count   int
}

// NewRollingStats returns a RollingStats.
func NewRollingStats(size int) *RollingStats {
	return &RollingStats{
		records: make([]float64, size),
		size:    size,
	}
}

// Add adds an element.
func (r *RollingStats) Add(n float64) {
	r.records[r.count%r.size] = n
	r.count++
}

// Median returns the median of the records.
// it can be used to filter noise.
// References: https://en.wikipedia.org/wiki/Median_filter.
func (r *RollingStats) Median() float64 {
	if r.count == 0 {
		return 0
	}
	records := r.records
	if r.count < r.size {
		records = r.records[:r.count]
	}
	median, _ := stats.Median(records)
	return median
}

// GetMean returns the mean of the records which have specified window size.
func GetMean(records []uint64, size uint32, count uint32) uint64 {
	if count == 0 {
		return 0
	}
	min := uint32(len(records))
	if size < min {
		min = size
	}
	if count < min {
		min = count
	}
	recordsRef := records[:min]

	var sum uint64
	for _, n := range recordsRef {
		sum += n
	}
	return sum / uint64(min)
}
