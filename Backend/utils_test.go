package main

import (
	"testing"
)

func Test_checkUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Empty username",
			args:    args{username: ""},
			wantErr: true,
		},
		{
			name:    "Too long username",
			args:    args{username: "johnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohns1"},
			wantErr: true,
		},
		{
			name:    "Too short username",
			args:    args{username: "ab"},
			wantErr: true,
		},
		{
			name:    "3 chars only",
			args:    args{username: "jon"},
			wantErr: false,
		},
		{
			name:    "chars with digits with underscore",
			args:    args{username: "john_doe1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkUsername(tt.args.username); (err != nil) != tt.wantErr {
				t.Errorf("checkUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Empty password",
			args:    args{password: ""},
			wantErr: true,
		},
		{
			name:    "Too short password",
			args:    args{password: "12345"},
			wantErr: true,
		},
		{
			name:    "Too long password",
			args:    args{password: "johnsjohnsjohnsjohnsjohnsjohns123"},
			wantErr: true,
		},
		{
			name:    "6 chars password",
			args:    args{password: "johns1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkPassword(tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("checkPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkPinCode(t *testing.T) {
	type args struct {
		pincode string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Empty pincode",
			args:    args{pincode: ""},
			wantErr: true,
		},
		{
			name:    "5 digits pincode",
			args:    args{pincode: "12345"},
			wantErr: true,
		},
		{
			name:    "Non digits pincode",
			args:    args{pincode: "azerty"},
			wantErr: true,
		},
		{
			name:    "6 digits pincode",
			args:    args{pincode: "123456"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkPinCode(tt.args.pincode); (err != nil) != tt.wantErr {
				t.Errorf("checkPinCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Bechmark_checkPinCode(b *testing.B) {
	type args struct {
		pincode string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Empty pincode",
			args:    args{pincode: ""},
			wantErr: true,
		},
		{
			name:    "5 digits pincode",
			args:    args{pincode: "12345"},
			wantErr: true,
		},
		{
			name:    "Non digits pincode",
			args:    args{pincode: "azerty"},
			wantErr: true,
		},
		{
			name:    "6 digits pincode",
			args:    args{pincode: "123456"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			checkPinCode(tt.args.pincode)
		})
	}
}

func Benchmark_checkUsername(b *testing.B) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Empty username",
			args:    args{username: ""},
			wantErr: true,
		},
		{
			name:    "Too long username",
			args:    args{username: "johnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohnsjohns1"},
			wantErr: true,
		},
		{
			name:    "Too short username",
			args:    args{username: "ab"},
			wantErr: true,
		},
		{
			name:    "3 chars only",
			args:    args{username: "jon"},
			wantErr: false,
		},
		{
			name:    "chars with digits with underscore",
			args:    args{username: "john_doe1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			checkUsername(tt.args.username)
		})
	}
}

func Benchmark_checkPassword(b *testing.B) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Empty password",
			args:    args{password: ""},
			wantErr: true,
		},
		{
			name:    "Too short password",
			args:    args{password: "12345"},
			wantErr: true,
		},
		{
			name:    "Too long password",
			args:    args{password: "johnsjohnsjohnsjohnsjohnsjohns123"},
			wantErr: true,
		},
		{
			name:    "6 chars password",
			args:    args{password: "johns1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			checkPassword(tt.args.password)
		})
	}
}
