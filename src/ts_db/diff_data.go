package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

import "github.com/sergi/go-diff/diffmatchpatch"

// DiffStore
// DiffStore.GetPrevious(timestamp)

func (self *DiffData) encode() ([]byte, error) {
	enc, err := json.Marshal(self)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (self *DiffData) decode(data []byte) error {
	err := json.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (self *DiffData) diffRebuildtexts(diffs []diffmatchpatch.Diff) []string {
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

func (self *DiffData) rebuildTextsToDiffN(n int) (string, error) {
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

func (self DiffData) GetImportantVersions() ([]versionsInfo, time.Duration) {
	m := map[int]int{}
	lastTime := time.Now().AddDate(0, -1, 0)
	totalTime := time.Now().Sub(time.Now())
	for i := range self.Diffs {
		parsedTime, _ := time.Parse(time.ANSIC, self.Timestamps[i])
		duration := parsedTime.Sub(lastTime)
		if duration.Minutes() < 3 {
			totalTime += duration
		}
		m[i] = int(duration.Seconds())
		if i > 0 {
			m[i-1] = m[i]
		}
		// On to the next one
		lastTime = parsedTime
	}

	// Sort in order of decreasing diff times
	n := map[int][]int{}
	var a []int
	for k, v := range m {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))

	// Get the top 4 biggest diff times
	var importantVersions []int
	var r []versionsInfo
	for _, k := range a {
		for _, s := range n[k] {
			if s != 0 && s != len(n) {
				// fmt.Printf("%d, %d\n", s, k)
				importantVersions = append(importantVersions, s)
				if len(importantVersions) > 10 {
					sort.Ints(importantVersions)
					for _, nn := range importantVersions {
						r = append(r, versionsInfo{self.Timestamps[nn], nn})
					}
					return r, totalTime
				}
			}
		}
	}
	sort.Ints(importantVersions)
	for _, nn := range importantVersions {
		r = append(r, versionsInfo{self.Timestamps[nn], nn})
	}
	return r, totalTime
}

func (self *DiffData) Update(newText string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(self.CurrentText, newText, true)
	delta := dmp.DiffToDelta(diffs)
	self.CurrentText = newText
	self.Timestamps = append(self.Timestamps, time.Now().Format(time.ANSIC))
	self.Diffs = append(self.Diffs, delta)
	self.Title = strings.ToLower(self.Title)
}

func (self *DiffData) GetCurrent() string {
	return self.CurrentText
}
