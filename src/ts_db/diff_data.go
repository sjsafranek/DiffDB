package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

import "github.com/sergi/go-diff/diffmatchpatch"

func (self *DiffStore) encode() ([]byte, error) {
	enc, err := json.Marshal(self)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (self *DiffStore) decode(data []byte) error {
	err := json.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

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

func (self *DiffStore) rebuildTextsToDiffN(timestamp int64, snapshots []int64) (string, error) {
	dmp := diffmatchpatch.New()
	lastText := ""
	self.lock.Lock()
	//for i, diff := range self.Diffs {
	for _, snapshot := range snapshots {

		diff := self.Diffs[snapshot]

		//log.Println(diff)

		seq1, _ := dmp.DiffFromDelta(lastText, diff)
		textsLinemode := self.diffRebuildtexts(seq1)
		rebuilt := textsLinemode[len(textsLinemode)-1]
		//if i == timestamp {
		if snapshot == timestamp {
			return rebuilt, nil
		}
		lastText = rebuilt
	}
	self.lock.Unlock()
	return "", fmt.Errorf("Could not rebuild from diffs")
}

func (self *DiffStore) Update(newText string) {
	self.lock.RLock()
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(self.CurrentValue, newText, true)
	delta := dmp.DiffToDelta(diffs)
	self.CurrentValue = newText
	now := time.Now().UnixNano()
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
	keys = MergeSort(keys)
	return keys
}

func (self *DiffStore) GetPreviousByTimestamp(timestamp int64) string {

	snapshots := self.GetSnapshots()

	// default to first value
	var ts int64 = snapshots[0]
	//var ts int64 = 0

	for _, snapshot := range snapshots {
		if timestamp >= snapshot && ts < snapshot {
			ts = snapshot
		}
	}

	oldValue, err := self.rebuildTextsToDiffN(ts, snapshots)
	if nil != err {
		log.Fatal(err)
	}

	return oldValue
}

func (self *DiffStore) GetPreviousByIndex(idx int) string {

	snapshots := self.GetSnapshots()

	if idx > len(snapshots)-1 {
		idx = len(snapshots) - 1
	}

	var ts int64 = snapshots[idx]

	oldValue, err := self.rebuildTextsToDiffN(ts, snapshots)
	if nil != err {
		log.Fatal(err)
	}

	return oldValue
}
