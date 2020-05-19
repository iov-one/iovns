package keeper

import (
	"reflect"
	"testing"
)

func Test_getDomainPrefixKey(t *testing.T) {
	type args struct {
		domainName string
	}
	tests := []struct {
		name        string
		args        args
		want        []byte
		expectPanic bool
	}{
		{
			name: "success",
			args: args{
				domainName: "test",
			},
			want: []byte{116, 101, 115, 116, 255},
		},
		{
			name: "panic",
			args: args{
				domainName: "test" + string(0xFF),
			},
			want:        nil,
			expectPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Fatalf("getDomainPrefixKey() panic expected")
					}
				}()
			}
			if got := getDomainPrefixKey(tt.args.domainName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDomainPrefixKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
