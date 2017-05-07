package diffstore

// MergeInt64 merges two arrays of int64 into a single array.
func MergeInt64(l, r []int64) []int64 {
	ret := make([]int64, 0, len(l)+len(r))
	for len(l) > 0 || len(r) > 0 {
		if len(l) == 0 {
			return append(ret, r...)
		}
		if len(r) == 0 {
			return append(ret, l...)
		}
		if l[0] <= r[0] {
			ret = append(ret, l[0])
			l = l[1:]
		} else {
			ret = append(ret, r[0])
			r = r[1:]
		}
	}
	return ret
}

// MergeSortInt64 merge sorts arrays of int64.
func MergeSortInt64(s []int64) []int64 {
	if len(s) <= 1 {
		return s
	}
	n := len(s) / 2
	l := MergeSortInt64(s[:n])
	r := MergeSortInt64(s[n:])
	return MergeInt64(l, r)
}
