package utils

import (
	"reflect"
	"testing"
)

type testType struct {
	Test string
}

func TestCloneType(t *testing.T) {
	type args struct {
		x interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "success",
			args: args{
				x: &testType{
					Test: "not empty",
				},
			},
			want: &testType{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CloneFromValue(tt.args.x); !reflect.DeepEqual(got.(*testType), tt.want) {
				t.Errorf("CloneFromValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCloneType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CloneFromValue(&testType{})
	}
}

func TestNewValueFromType(t *testing.T) {
	x := &testType{}
	type args struct {
		typ reflect.Type
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "success",
			args: args{
				typ: reflect.ValueOf(x).Type().Elem(),
			},
			want: &testType{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CloneFromType(tt.args.typ); !reflect.DeepEqual(got.(*testType), tt.want) {
				t.Errorf("CloneFromType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNewValueFromType(b *testing.B) {
	typ := reflect.ValueOf(&testType{}).Type().Elem()
	//
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CloneFromType(typ)
	}
}
