package helper

import (
	"testing"
)

// go test -run Add -v  指定运行的方法，是正则，这样会运行 *Add*
// 获取测试覆盖率 go test -v -coverprofile cover.out     | go tool cover -html cover.out -o cover.html

func TestAdd2(t *testing.T) {
	tests := []struct {
		param1 int
		param2 int
		result int
	}{
		{
			param1: 1,
			param2: 1,
			result: 2,
		},
		{
			param1: 1,
			param2: 1,
			result: 2,
		},
		{
			param1: 1,
			param2: 1,
			result: 2,
		},
	}

	for _, test := range tests {
		r := Add2(test.param1, test.param2)
		if r != test.result {
			t.Errorf("error: expecte %d, but got %d", test.result, r)
		}
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		param1 int
		param2 int
		result int
	}{
		{
			param1: 1,
			param2: 1,
			result: 2,
		},
		{
			param1: 1,
			param2: 1,
			result: 2,
		},
		{
			param1: 1,
			param2: 1,
			result: 2,
		},
	}

	for _, test := range tests {
		r := Add(test.param1, test.param2)
		if r != test.result {
			t.Errorf("error: expecte %d, but got %d", test.result, r)
		}
	}
}
