package simpleapi

import (
	"github.com/redagain/go-testrequest"
	"net/http"
	"reflect"
	"testing"
)

func TestBookRequest_IsValid(t *testing.T) {
	type fields struct {
		Title   string
		ISBN    string
		Authors []string
	}
	tests := []struct {
		name   string
		fields fields
		wantOk bool
	}{
		{
			name: "EmptyTitle",
			fields: fields{
				Title: "",
			},
			wantOk: false,
		},
		{
			name: "EmptyISBN",
			fields: fields{
				Title: "The Go Programming Language",
				ISBN:  "",
			},
			wantOk: false,
		},
		{
			name: "EmptyAuthors",
			fields: fields{
				Title:   "The Go Programming Language",
				ISBN:    "978-0134190440",
				Authors: []string{},
			},
			wantOk: false,
		},
		{
			name: "EmptyAuthorValues",
			fields: fields{
				Title:   "The Go Programming Language",
				ISBN:    "978-0134190440",
				Authors: []string{""},
			},
			wantOk: false,
		},
		{
			name: "Ok",
			fields: fields{
				Title:   "The Go Programming Language",
				ISBN:    "978-0134190440",
				Authors: []string{"Alan A. A. Donovan", "Brian W. Kernighan"},
			},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &BookRequest{
				Title:   tt.fields.Title,
				ISBN:    tt.fields.ISBN,
				Authors: tt.fields.Authors,
			}
			if gotOk := req.IsValid(); gotOk != tt.wantOk {
				t.Errorf("IsValid() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestNewBookRequest(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *BookRequest
		wantErr bool
	}{
		{
			name: "FromEmptyBody",
			args: args{
				req: testrequest.Builder().Request(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "FromPostForm",
			args: args{
				req: testrequest.Builder().
					SetPostFormValue("title", "The Go Programming Language").
					SetPostFormValue("isbn", "978-0134190440").
					SetPostFormValue("authors", "Alan A. A. Donovan", "Brian W. Kernighan").
					Request(),
			},
			want: &BookRequest{
				Title:   "The Go Programming Language",
				ISBN:    "978-0134190440",
				Authors: []string{"Alan A. A. Donovan", "Brian W. Kernighan"},
			},
			wantErr: false,
		},
		{
			name: "FromInvalidJSON",
			args: args{
				req: testrequest.Builder().
					SetJSON([]byte("{")).
					Request(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "FromJSON",
			args: args{
				req: testrequest.Builder().SetJSONFromValue(
					map[string]interface{}{
						"title":   "The Go Programming Language",
						"isbn":    "978-0134190440",
						"authors": []string{"Alan A. A. Donovan", "Brian W. Kernighan"},
					}).Request(),
			},
			want: &BookRequest{
				Title:   "The Go Programming Language",
				ISBN:    "978-0134190440",
				Authors: []string{"Alan A. A. Donovan", "Brian W. Kernighan"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBookRequest(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBookFromRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateBookFromRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
