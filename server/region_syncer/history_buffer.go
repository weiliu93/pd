// copyright 2018 pingcap, inc.
//
// licensed under the apache license, version 2.0 (the "license");
// you may not use this file except in compliance with the license.
// you may obtain a copy of the license at
//
//     http://www.apache.org/licenses/license-2.0
//
// unless required by applicable law or agreed to in writing, software
// distributed under the license is distributed on an "as is" basis,
// see the license for the specific language governing permissions and
// limitations under the license.

package syncer

import (
	"strconv"

	"github.com/pingcap/pd/server/core"
	log "github.com/sirupsen/logrus"
)

const (
	historyKey        = "historyIndex"
	defaultFlushCount = 100
)

type historyBuffer struct {
	index      uint64
	records    []*core.RegionInfo
	head       int
	tail       int
	size       int
	kv         core.KVBase
	flushCount int
}

func newHistoryBuffer(size int, kv core.KVBase) *historyBuffer {
	// use an empty space to simplify operation
	size++
	if size < 2 {
		size = 2
	}
	records := make([]*core.RegionInfo, size)
	h := &historyBuffer{
		records: records,
		size:    size,
		kv:      kv,
	}
	h.reload()
	return h
}

func (h *historyBuffer) len() int {
	if h.tail < h.head {
		return h.tail + h.size - h.head
	}
	return h.tail - h.head
}

func (h *historyBuffer) lastIndex() uint64 {
	return h.index
}

func (h *historyBuffer) firstIndex() uint64 {
	return h.index - uint64(h.len())
}

func (h *historyBuffer) record(r *core.RegionInfo) {
	if (h.tail+1)%h.size == h.head {
		h.head = (h.head + 1) % h.size
	}
	h.tail = (h.tail + 1) % h.size
	h.records[h.tail] = r
	h.index++
	h.flushCount--
	if h.flushCount <= 0 {
		h.persist()
	}
}

func (h *historyBuffer) get(index uint64) *core.RegionInfo {
	if index <= h.lastIndex() && index > h.firstIndex() {
		return h.records[index%uint64(h.size)]
	}
	return nil
}

func (h *historyBuffer) reload() {
	v, err := h.kv.Load(historyKey)
	if err != nil {
		log.Warnf("load history index failed: %s", err)
	}
	h.index, err = strconv.ParseUint(v, 10, 64)
	if err != nil {
		log.Warnf("load history index failed: %s", err)
	}
	h.flushCount = defaultFlushCount
}

func (h *historyBuffer) persist() {
	err := h.kv.Save(historyKey, strconv.FormatUint(h.lastIndex(), 10))
	if err != nil {
		log.Warnf("persist history index (%d) failed: %v", h.lastIndex(), err)
	}
}
