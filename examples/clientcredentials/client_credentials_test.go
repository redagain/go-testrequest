package clientcredentials

import (
	"github.com/redagain/go-testrequest"
	"net/http"
	"testing"
)

func TestBasicCredentials(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name         string
		args         args
		wantClientID string
		wantSecret   string
		wantOk       bool
	}{
		{
			name: "EmptyRequest",
			args: args{
				req: testrequest.Builder().Request(),
			},
			wantClientID: "",
			wantSecret:   "",
			wantOk:       false,
		},
		{
			name: "EmptyCredentialsFromBasicAuth",
			args: args{
				req: testrequest.Builder().SetBasicAuth("", "").Request(),
			},
			wantClientID: "",
			wantSecret:   "",
			wantOk:       true,
		},
		{
			name: "CredentialsFromBasicAuth",
			args: args{
				req: testrequest.Builder().SetBasicAuth("test", "test").Request(),
			},
			wantClientID: "test",
			wantSecret:   "test",
			wantOk:       true,
		},
		{
			name: "NotContainSecret",
			args: args{
				req: testrequest.Builder().
					SetPostFormValue("test", "test").
					Request(),
			},
			wantClientID: "",
			wantSecret:   "",
			wantOk:       false,
		},
		{
			name: "EmptyCredentialsFromForm",
			args: args{
				req: testrequest.Builder().
					SetPostFormValue("client_id", "").
					SetPostFormValue("client_secret", "").
					Request(),
			},
			wantClientID: "",
			wantSecret:   "",
			wantOk:       true,
		},
		{
			name: "CredentialsFromForm",
			args: args{
				req: testrequest.Builder().
					SetQueryValue("client_id", "test").
					SetPostFormValue("client_secret", "test").
					Request(),
			},
			wantClientID: "test",
			wantSecret:   "test",
			wantOk:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClientID, gotSecret, gotOk := BasicCredentials(tt.args.req)
			if gotClientID != tt.wantClientID {
				t.Errorf("BasicCredentials() gotClientID = %v, want %v", gotClientID, tt.wantClientID)
			}
			if gotSecret != tt.wantSecret {
				t.Errorf("BasicCredentials() gotSecret = %v, want %v", gotSecret, tt.wantSecret)
			}
			if gotOk != tt.wantOk {
				t.Errorf("BasicCredentials() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
