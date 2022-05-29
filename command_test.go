package gshell

import (
	"reflect"
	"testing"
)

func Test_samplize(t *testing.T) {
	type TestStructChild struct {
		ChildInt32  int32
		ChildString string
	}
	type TestStruct struct {
		String    string
		Int32     int32
		ChildList []*TestStructChild
		Child     *TestStructChild
		Int32List []int32
	}

	r := &TestStruct{}
	samplize(r)

	want := &TestStruct{
		String:    "hoge",
		Int32:     1,
		ChildList: []*TestStructChild{{ChildInt32: 1, ChildString: "hoge"}},
		Child:     &TestStructChild{ChildInt32: 1, ChildString: "hoge"},
		Int32List: []int32{1},
	}
	if !reflect.DeepEqual(r, want) {
		t.Errorf("TestStruct got = %v, want %v", r, want)
	}
}
