package gshell

import (
	"reflect"
	"testing"
)

func Test_samplize(t *testing.T) {
	type AliasInt int
	type TestStructChild struct {
		ChildInt32  int32
		ChildString string
	}
	type TestStruct struct {
		String    string
		Int32     int32
		ChildList []*TestStructChild
		Child     *TestStructChild
		Int64List []int64
		AliasInt  []AliasInt
	}

	type args struct {
		r interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "正常変換",
			args: args{&TestStruct{}},
			want: &TestStruct{
				String:    "hoge",
				Int32:     1,
				ChildList: []*TestStructChild{{ChildInt32: 1, ChildString: "hoge"}},
				Child:     &TestStructChild{ChildInt32: 1, ChildString: "hoge"},
				Int64List: []int64{1},
				AliasInt:  []AliasInt{1},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			samplize(tt.args.r)
		})
		if !reflect.DeepEqual(tt.args.r, tt.want) {
			t.Errorf("TestStruct got = %v, want %v", tt.args.r, tt.want)
		}
	}
}
