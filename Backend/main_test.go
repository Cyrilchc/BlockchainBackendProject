package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_createUserAndWallet(t *testing.T) {
	httpPostUrl := fmt.Sprintf("http://%s:%s/", HOST, SERVER_PORT)
	httpParams, _ := json.Marshal(
		Player{
			Username: "johndoe",
			Password: "Foobar.1",
			Pincode:  "123456",
		})
	httpRequest, _ := http.NewRequest(http.MethodPost, httpPostUrl, bytes.NewBuffer(httpParams))
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Normal query",
			args: args{
				w: httptest.NewRecorder(),
				r: httpRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createUserAndWallet(tt.args.w, tt.args.r)
		})
	}
}
