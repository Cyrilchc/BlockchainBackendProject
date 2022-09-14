package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func Test_buildHttpErrorMessage(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Empty message",
			args:    args{message: ""},
			wantErr: false,
			want:    "{\"error\":\"\"}",
		},
		{
			name:    "With message",
			args:    args{message: "Unable to insert player"},
			wantErr: false,
			want:    "{\"error\":\"Unable to insert player\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildHttpErrorMessage(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildHttpErrorMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("buildHttpErrorMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkHttpMethod(t *testing.T) {
	type args struct {
		method string
		w      http.ResponseWriter
		r      *http.Request
	}
	var tests = []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "With incorrect Method",
			args: args{
				method: "PATCH",
				w:      httptest.NewRecorder(),
				r:      httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%s/", HOST, SERVER_PORT), bytes.NewBuffer([]byte{})),
			},
			wantErr: true,
		},
		{
			name: "With correct Method",
			args: args{
				method: "POST",
				w:      httptest.NewRecorder(),
				r:      httptest.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/", HOST, SERVER_PORT), bytes.NewBuffer([]byte{})),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkHttpMethod(tt.args.method, tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("checkHttpMethod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Benchmark_checkHttpMethod(b *testing.B) {
	type args struct {
		method string
		w      http.ResponseWriter
		r      *http.Request
	}
	var tests = []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "With incorrect Method",
			args: args{
				method: "PATCH",
				w:      httptest.NewRecorder(),
				r:      httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%s/", HOST, SERVER_PORT), bytes.NewBuffer([]byte{})),
			},
			wantErr: true,
		},
		{
			name: "With correct Method",
			args: args{
				method: "POST",
				w:      httptest.NewRecorder(),
				r:      httptest.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/", HOST, SERVER_PORT), bytes.NewBuffer([]byte{})),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			checkHttpMethod(tt.args.method, tt.args.w, tt.args.r)
		})
	}
}

func Benchmark_buildHttpErrorMessage(b *testing.B) {
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Empty message",
			args:    args{message: ""},
			wantErr: false,
			want:    "{\"error\":\"\"}",
		},
		{
			name:    "With message",
			args:    args{message: "Unable to insert player"},
			wantErr: false,
			want:    "{\"error\":\"Unable to insert player\"}",
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			buildHttpErrorMessage(tt.args.message)
		})
	}
}

func Benchmark_checkPinCode(b *testing.B) {
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
