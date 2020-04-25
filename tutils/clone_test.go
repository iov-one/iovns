package tutils

import (
	"github.com/iov-one/iovns/x/domain/types"
	"reflect"
	"testing"
)

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
				x: &types.Account{
					Domain: "test",
				},
			},
			want: &types.Account{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CloneFromValue(tt.args.x); !reflect.DeepEqual(got.(*types.Account), tt.want) {
				t.Errorf("CloneFromValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCloneType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CloneFromValue(&types.Account{})
	}
}

func TestNewValueFromType(t *testing.T) {
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
				typ: reflect.ValueOf(&types.Account{Name: "test"}).Type().Elem(),
			},
			want: &types.Account{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CloneFromType(tt.args.typ); !reflect.DeepEqual(got.(*types.Account), tt.want) {
				t.Errorf("CloneFromType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNewValueFromType(b *testing.B) {
	typ := reflect.ValueOf(&types.Account{}).Type().Elem()
	//
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CloneFromType(typ)
	}
}
