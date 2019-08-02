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
	"testing"

	. "github.com/pingcap/check"
)

func Test(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&testRollingStats{})

type testRollingStats struct{}

func (t *testRollingStats) TestRollingMedian(c *C) {
	data := []float64{2, 4, 2, 800, 600, 6, 3}
	expected := []float64{2, 3, 2, 3, 4, 6, 6}
	stats := NewRollingStats(5)
	for i, e := range data {
		stats.Add(e)
		c.Assert(stats.Median(), Equals, expected[i])
	}
}

func (t *testRollingStats) TestMedianAndMean(c *C) {
	data := []uint64{2, 4, 2, 8, 7, 6, 0}

	mean := GetMean(data, uint32(len(data)), uint32(len(data)))
	c.Assert(mean, Equals, 4)
}
