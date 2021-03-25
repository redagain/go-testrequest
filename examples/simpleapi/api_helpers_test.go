package simpleapi

import (
	"github.com/redagain/go-testrequest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBearerAuth(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name            string
		args            args
		wantAccessToken string
		wantOk          bool
	}{
		{
			name:            "ShouldBeFalseIfAuthorizationHeaderIsMissing",
			args:            args{req: testrequest.Builder().Request()},
			wantAccessToken: "",
			wantOk:          false,
		},
		{
			name:            "ShouldBeFalseIfAuthorizationHeaderIsEmpty",
			args:            args{req: testrequest.Builder().SetBearerAuth("").Request()},
			wantAccessToken: "",
			wantOk:          false,
		},
		{
			name: "ShouldBeNonEmptyAccessTokenAndOk",
			args: args{
				req: testrequest.Builder().SetBearerAuth("CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ").Request(),
			},
			wantAccessToken: "CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ",
			wantOk:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAccessToken, gotOk := BearerAuth(tt.args.req)
			if gotAccessToken != tt.wantAccessToken {
				t.Errorf("BearerAuth() gotAccessToken = %v, want %v", gotAccessToken, tt.wantAccessToken)
			}
			if gotOk != tt.wantOk {
				t.Errorf("BearerAuth() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestCreated(t *testing.T) {
	type args struct {
		w    *httptest.ResponseRecorder
		resp interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantStatus  int
		wantContent string
		wantErr     bool
	}{
		{
			name: "NilResponse",
			args: args{
				w:    httptest.NewRecorder(),
				resp: nil,
			},
			wantErr:     false,
			wantStatus:  http.StatusCreated,
			wantContent: "application/json;charset=UTF-8",
		},
		{
			name: "GoodResponse",
			args: args{
				w:    httptest.NewRecorder(),
				resp: map[string]interface{}{"id": 1},
			},
			wantErr:     false,
			wantStatus:  http.StatusCreated,
			wantContent: "application/json;charset=UTF-8",
		},
		{
			name: "InvalidResponse",
			args: args{
				w:    httptest.NewRecorder(),
				resp: func() {},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Created(tt.args.w, tt.args.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Created() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				resp := tt.args.w
				gotStatus := resp.Code
				if gotStatus != tt.wantStatus {
					t.Errorf("Created() status = %v, wantStatus %v", gotStatus, tt.wantStatus)
					return
				}
				gotContent := resp.Header().Get("Content-Type")
				if gotContent != tt.wantContent {
					t.Errorf("Created() content = %v, wantContent %v", gotStatus, tt.wantStatus)
					return
				}
			}
		})
	}
}
