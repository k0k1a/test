package src

import (
	"reflect"
	"testing"
)

type TestCase struct {
	Name     string
	Input    []string
	Expected []bool
	Clear    bool
}

func TestCacheManager(t *testing.T) {
	var testcases1 = []TestCase{
		{
			Name:     "all same keys",
			Input:    []string{"a", "a", "a"},
			Expected: []bool{false, true, true},
		},
		{
			Name:     "all different keys",
			Input:    []string{"a", "b", "c"},
			Expected: []bool{true, false, false},
		},
		{
			Name:     "clear all keys",
			Input:    []string{"a", "b", "c"},
			Expected: []bool{false, false, false},
			Clear:    true,
		},
		{
			Name:     "add same key",
			Input:    []string{"a", "b", "c"},
			Expected: []bool{true, true, true},
		},
	}
	doTest(t, testcases1)

	var testcases2 = []TestCase{
		{
			Name:     "empty",
			Input:    []string{},
			Expected: []bool{},
		},
		{
			Name:     "all different keys",
			Input:    []string{"a", "b", "c"},
			Expected: []bool{false, false, false},
		},
	}
	doTest(t, testcases2)
}

func doTest(t *testing.T, testcases []TestCase) {
	server, err := StartHttpsServer(":8848")
	defer server.server.Close()
	if err != nil {
		t.Fatalf("start server err %v", err)
	}
	for _, item := range testcases {
		t.Run(item.Name, func(t *testing.T) {
			if item.Clear {
				server.cacheManager.Clear()
			}

			ret, err := BcjClient(item.Input)
			if err != nil {
				t.Fatalf("test case %v err %v", item.Name, err)
			}

			if !reflect.DeepEqual(ret, item.Expected) {
				t.Fatalf("test case %v expected %v got %v", item.Name, item.Expected, ret)
			}
		})
	}
}
