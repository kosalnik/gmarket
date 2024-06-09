package crypt

import (
	"testing"
)

func TestPasswordHasher_Hash(t *testing.T) {
	type args struct {
		pwd string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{pwd: "123"},
			want: "202cb962ac59075b964b07152d234b70",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PasswordHasher{}
			got := p.Hash(tt.args.pwd)
			if got != tt.want {
				t.Errorf("Hash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordHasher_IsEquals(t *testing.T) {
	type args struct {
		pwd string
		h   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success",
			args: args{pwd: "123", h: "202cb962ac59075b964b07152d234b70"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPasswordHasher()
			if got := p.IsEquals(tt.args.pwd, tt.args.h); got != tt.want {
				t.Errorf("IsEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}
