package diff_store

import (
	"encoding/json"
	"fmt"
	"time"
)

import "github.com/sergi/go-diff/diffmatchpatch"

// NewDiffStore creates and returns DiffStore struct
func NewDiffStore(name string) DiffStore {
	var ddata DiffStore
	ddata.Name = name
	ddata.CurrentValue = ""
	ddata.Diffs = make(map[int64]string)
	return ddata
}

// Encode marshals struct into json.
func (self *DiffStore) Encode() ([]byte, error) {
	enc, err := json.Marshal(self)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// Decode unmarshals struct from json.
func (self *DiffStore) Decode(data []byte) error {
	err := json.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

// diffRebuildtexts rebuilds text value from list of diff changes.
func (self *DiffStore) diffRebuildtexts(diffs []diffmatchpatch.Diff) []string {
	text := []string{"", ""}
	for _, diff := range diffs {
		if diff.Type != diffmatchpatch.DiffInsert {
			text[0] += diff.Text
		}
		if diff.Type != diffmatchpatch.DiffDelete {
			text[1] += diff.Text
		}
	}
	return text
}

// rebuildTextsToDiffN rebuilds text value from diff changes until timestamp is reached.
func (self *DiffStore) rebuildTextsToDiffN(timestamp int64, snapshots []int64) (string, error) {
	dmp := diffmatchpatch.New()
	lastText := ""
	self.lock.Lock()

	for _, snapshot := range snapshots {

		diff := self.Diffs[snapshot]
		seq1, _ := dmp.DiffFromDelta(lastText, diff)
		textsLinemode := self.diffRebuildtexts(seq1)
		rebuilt := textsLinemode[len(textsLinemode)-1]

		if snapshot == timestamp {
			self.lock.Unlock()
			return rebuilt, nil
		}
		lastText = rebuilt
	}

	self.lock.Unlock()
	return "", fmt.Errorf("Could not rebuild from diffs")
}

// Update adds new text change.
func (self *DiffStore) Update(newText string) {

	// check for changes
	if self.GetCurrent() == newText {
		return
	}

	self.lock.RLock()
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(self.CurrentValue, newText, true)
	delta := dmp.DiffToDelta(diffs)
	self.CurrentValue = newText
	now := time.Now().UnixNano()

	if nil == self.Diffs {
		self.Diffs = make(map[int64]string)
	}

	self.Diffs[now] = delta
	self.lock.RUnlock()
}

// GetCurrent returns current text value.
func (self *DiffStore) GetCurrent() string {
	return self.CurrentValue
}

// GetSnapshots returns a list of UnixNano timestamps for snapshots.
func (self *DiffStore) GetSnapshots() []int64 {
	self.lock.Lock()
	keys := make([]int64, 0, len(self.Diffs))
	for k := range self.Diffs {
		keys = append(keys, k)
	}
	self.lock.Unlock()
	// SORT KEYS
	keys = MergeSortInt64(keys)
	return keys
}

// GetPreviousByTimestamp returns text value at given timestamp.
func (self *DiffStore) GetPreviousByTimestamp(timestamp int64) (string, error) {

	// check inputs
	if 0 > timestamp {
		return "", fmt.Errorf("Timestamps most be positive integer")
	}

	// get change snapshot
	snapshots := self.GetSnapshots()

	// default to first value
	var ts int64 = snapshots[0]

	// find timestamp
	for _, snapshot := range snapshots {
		if timestamp >= snapshot && ts < snapshot {
			ts = snapshot
		}
	}

	// use timestamp to find value
	oldValue, err := self.rebuildTextsToDiffN(ts, snapshots)
	return oldValue, err
}

// GetPreviousByIndex returns value at given index.
func (self *DiffStore) GetPreviousByIndex(idx int) (string, error) {

	// check inputs
	if 0 > idx {
		return "", fmt.Errorf("Index most be positive integer")
	}

	// get change snapshots
	snapshots := self.GetSnapshots()

	// if index greater than length of snapshot
	// default to last snapshot
	if idx > len(snapshots)-1 {
		idx = len(snapshots) - 1
	}

	// use index to find timestamp
	var ts int64 = snapshots[idx]

	// use timestamp to find value
	oldValue, err := self.rebuildTextsToDiffN(ts, snapshots)
	return oldValue, err
}

// GetPreviousWithinRange returns text values within a given timestamp range
func (self *DiffStore) GetPreviousWithinTimestampRange(begin_timestamp int64, end_timestamp int64) (map[int64]string, error) {

	// TODO:
	// - Calculate old values i one pass

	values := make(map[int64]string)

	if begin_timestamp > end_timestamp {
		return values, fmt.Errorf("begin_timestamp must be greater than end_timestamp")
	}

	// check inputs
	if 0 > begin_timestamp || 0 > end_timestamp {
		return values, fmt.Errorf("Timestamps most be positive integers")
	}

	// rebuild all values within range
	snapshots := self.GetSnapshots()
	for _, snapshot := range snapshots {
		if begin_timestamp <= snapshot && end_timestamp >= snapshot {
			value, err := self.rebuildTextsToDiffN(snapshot, snapshots)
			if nil != err {
				return values, err
			}

			values[snapshot] = value
		}
	}

	// return values
	return values, nil
}
