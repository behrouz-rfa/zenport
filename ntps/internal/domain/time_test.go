package domain

import (
	"reflect"
	"testing"
)

func TestCreateTime(t *testing.T) {
	type args struct {
		id    string
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    *Time
		wantErr bool
	}{
		{
			name: "Ask Time Correctly",
			args: args{
				id:    "baea7756-f814-4d34-b3b9-7549603dc721",
				input: "What time is it?",
			},
			want:    NewNtp("baea7756-f814-4d34-b3b9-7549603dc721"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTime(tt.args.id, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ID(), tt.want.ID()) {
				t.Errorf("CreateTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}
