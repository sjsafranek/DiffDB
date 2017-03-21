package skeleton_db

import (
	//"errors"
	"log"
	"testing"
)

// go test -bench=.
// go test -bench=. -test.benchmem

var testDiffDb DiffDb

func init() {
	testDiffDb = NewDiffDb("test.db")
}

func BenchmarkDbCreateKey(b *testing.B) {
	key, _ := NewUUID2()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testDiffDb.Load(key)
		if nil != err {
			log.Fatal(err)
		}
	}
}

func BenchmarkDbLoadKey(b *testing.B) {
	key, _ := NewUUID2()

	ddata, err := testDiffDb.Load(key)
	if nil != err {
		log.Fatal(err)
	}
	ddata.Update("test")
	testDiffDb.Save(ddata)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testDiffDb.Load(key)
		if nil != err {
			log.Fatal(err)
		}
	}
}

func BenchmarkDbUpdateKey(b *testing.B) {
	key, _ := NewUUID2()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ddata, err := testDiffDb.Load(key)
		if nil != err {
			log.Fatal(err)
		}
		ddata.Update("test")
	}
}

func BenchmarkDbUpdateAndSaveKey(b *testing.B) {
	key, _ := NewUUID2()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ddata, err := testDiffDb.Load(key)
		if nil != err {
			log.Fatal(err)
		}
		ddata.Update("test")
		testDiffDb.Save(ddata)
	}
}

func BenchmarkDbGetCurrentValueFromKey(b *testing.B) {
	key, _ := NewUUID2()

	ddata, err := testDiffDb.Load(key)
	if nil != err {
		log.Fatal(err)
	}
	ddata.Update("test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ddata.GetCurrent()
	}
}

func BenchmarkDbListKeySnapshots(b *testing.B) {
	key, _ := NewUUID2()

	ddata, err := testDiffDb.Load(key)
	if nil != err {
		log.Fatal(err)
	}
	ddata.Update("test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ddata.GetSnapshots()
	}
}

func BenchmarkDbPreviousSnapshot(b *testing.B) {
	key, _ := NewUUID2()

	ddata, err := testDiffDb.Load(key)
	if nil != err {
		log.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		newVal, _ := NewUUID2()
		ddata.Update(newVal)
		testDiffDb.Save(ddata)
	}

	snapshots := ddata.GetSnapshots()
	l := int(len(snapshots) / 2)
	timestamp := snapshots[l]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ddata, err := testDiffDb.Load(key)
		if nil != err {
			log.Fatal(err)
		}
		ddata.GetPreviousByTimestamp(timestamp)
	}

}
