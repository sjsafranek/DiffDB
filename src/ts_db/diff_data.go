package main

import (
	"encoding/json"
	"fmt"
	"log"
	//"sort"
	//"strings"
	"time"
)

import "github.com/sergi/go-diff/diffmatchpatch"

// DiffStore
// DiffStore.GetPrevious(timestamp)

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
	for _, myDiff := range diffs {
		if myDiff.Type != diffmatchpatch.DiffInsert {
			text[0] += myDiff.Text
		}
		if myDiff.Type != diffmatchpatch.DiffDelete {
			text[1] += myDiff.Text
		}
	}
	return text
}

func (self *DiffStore) rebuildTextsToDiffN(n int64) (string, error) {
	dmp := diffmatchpatch.New()
	lastText := ""
	for i, diff := range self.Diffs {
		seq1, _ := dmp.DiffFromDelta(lastText, diff)
		textsLinemode := self.diffRebuildtexts(seq1)
		rebuilt := textsLinemode[len(textsLinemode)-1]
		if i == n {
			return rebuilt, nil
		}
		lastText = rebuilt
	}
	return "", fmt.Errorf("Could not rebuild from diffs")
}

func (self *DiffStore) Update(newText string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(self.CurrentText, newText, true)
	delta := dmp.DiffToDelta(diffs)
	self.CurrentText = newText
	now := time.Now().UnixNano()
	self.Diffs[now] = delta
	//self.Title = strings.ToLower(self.Title)
}

func (self *DiffStore) GetCurrent() string {
	return self.CurrentText
}

func (self *DiffStore) GetSnapshots() []int64 {
	keys := make([]int64, 0, len(self.Diffs))
	for k := range self.Diffs {
		keys = append(keys, k)
	}
	return keys
}

func (self *DiffStore) GetPrevious(timestamp int64) string {
	var ts int64 = 0
	for i := range self.Diffs {
		if timestamp >= i && ts < i {
			ts = i
		}
	}
	oldValue, err := self.rebuildTextsToDiffN(ts)
	if nil != err {
		log.Fatal(err)
	}
	return oldValue
}
