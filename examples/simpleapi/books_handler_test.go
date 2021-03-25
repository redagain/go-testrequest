package simpleapi

import (
	"errors"
	"github.com/redagain/go-testrequest"
	"net/http"
	"testing"
)

func TestBooksHandler_Post(t *testing.T) {
	type fields struct {
		addBook func(req *BookRequest) (id string, err error)
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "FailedBookRequestCreation",
			fields: fields{},
			args: args{
				r: testrequest.Builder().Request(),
			},
			wantErr: true,
		},
		{
			name:   "InvalidBookRequest",
			fields: fields{},
			args: args{
				r: testrequest.Builder().SetJSONFromValue(
					map[string]interface{}{
						"title": "",
					}).Request(),
			},
			wantErr: true,
		},
		{
			name: "FailedAddBookInStore",
			fields: fields{
				addBook: func(req *BookRequest) (id string, err error) {
					err = errors.New("failed to create book in store")
					return
				},
			},
			args: args{
				r: testrequest.Builder().SetJSONFromValue(
					map[string]interface{}{
						"title":   "The Go Programming Language",
						"isbn":    "978-0134190440",
						"authors": []string{"Alan A. A. Donovan", "Brian W. Kernighan"},
					}).Request(),
			},
			wantErr: true,
		},
		{
			name: "Success",
			fields: fields{
				addBook: func(req *BookRequest) (id string, err error) {
					id = "1"
					return
				},
			},
			args: args{
				w: testrequest.NopResponseWriter(),
				r: testrequest.Builder().SetJSONFromValue(
					map[string]interface{}{
						"title":   "The Go Programming Language",
						"isbn":    "978-0134190440",
						"authors": []string{"Alan A. A. Donovan", "Brian W. Kernighan"},
					}).Request(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BooksHandler{
				addBook: tt.fields.addBook,
			}
			if err := h.Post(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
