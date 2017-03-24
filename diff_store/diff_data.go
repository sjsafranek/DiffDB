package diff_store

import (
	"encoding/json"
	"fmt"
	"time"
)

import "github.com/sergi/go-diff/diffmatchpatch"

// @method 		NewDiffStore
// @description Creates and returns DiffStore structs
// @returns 	DiffStore
func NewDiffStore(name string) DiffStore {
	var ddata DiffStore
	ddata.Name = name
	ddata.CurrentValue = ""
	ddata.Diffs = make(map[int64]string)
	return ddata
}

// @method 		encode
// @description encodes struct to json []byte
// @returns 	[]byte, error
func (self *DiffStore) Encode() ([]byte, error) {
	enc, err := json.Marshal(self)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// @method 		decode
// @description decodes json []btye to struct
// @param		{[]byte} data to decode
// @returns 	error
func (self *DiffStore) Decode(data []byte) error {
	err := json.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

// @method 		diffRebuildtexts
// @description Builds text value from changes
// @param		{[]diffmatchpatch.Diff} list of diff changes
// @returns 	string, error
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

// @method 		rebuildTextsToDiffN
// @description Builds text value from changes
// @param		{int64} unix nano timestamp
// @param		{[]int64} list of unix nano timestamps
// @returns 	string, error
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

// @method 		Update
// @description Updates value
// @param		string
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

// @method 		GetCurrent
// @description Returns current value
// @return 		string
func (self *DiffStore) GetCurrent() string {
	return self.CurrentValue
}

// @method 		GetSnapshots
// @description Returns a list of UnixNano timestamps for snapshots
// @return 		[]int64
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

// @method 		GetPreviousByTimestamp
// @description Returns value at given timestamp
// @param		{int64}
// @return 		string
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

// @method 		GetPreviousByIndex
// @description Returns value at given index
// @param		{int}
// @return 		string
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

// @method 		GetPreviousWithinRange
// @description Returns value at given timestamp
// @param		{int64} begin_timestamp
// @param		{int64} end_timestamp
// @return 		string
func (self *DiffStore) GetPreviousWithinTimestampRange(begin_timestamp int64, end_timestamp int64) (map[int64]string, error) {

	// TODO:
	// - Calculate old values i one pass

	values := make(map[int64]string)

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
