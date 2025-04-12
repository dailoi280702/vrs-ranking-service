package converter

import (
	"strconv"
)

func StringToInt64Slice(strings []string) ([]int64, error) {
	int64s := make([]int64, len(strings))
	for i, s := range strings {
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		int64s[i] = val
	}
	return int64s, nil
}

func StringToInt64SliceIgnoreError(strings []string) []int64 {
	var int64s []int64
	for _, s := range strings {
		val, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			int64s = append(int64s, val)
		}
	}
	return int64s
}
