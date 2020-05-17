package index

import (
	"reflect"
	"testing"
)

func TestPackUnpackBytes(t *testing.T) {
	type args struct {
		keys [][]byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				keys: [][]byte{[]byte("test"), []byte("another-key"), []byte("another-one-yet")},
			},
			want:    []byte{uint8(4), 116, 101, 115, 116, 11, 97, 110, 111, 116, 104, 101, 114, 45, 107, 101, 121, 15, 97, 110, 111, 116, 104, 101, 114, 45, 111, 110, 101, 45, 121, 101, 116},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PackBytes(tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("PackBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PackBytes() got = %v, want %v", got, tt.want)
			}
			// convert
			resp, err := UnpackBytes(got)
			if err != nil {
				t.Errorf("UnpackBytes() error = %v", err)
			}
			if !reflect.DeepEqual(resp, tt.args.keys) {
				t.Fatalf("UnpackBytes() wanted %v, got %v", tt.args.keys, resp)
			}
		})
	}
}
