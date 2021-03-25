package simpleapi

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrorStatusCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "Unauthorized",
			args: args{
				err: ErrUnauthorized,
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "BadRequest",
			args: args{
				err: ErrBadRequest,
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "UnsupportedMediaType",
			args: args{
				err: ErrUnsupportedMediaType,
			},
			wantStatusCode: http.StatusUnsupportedMediaType,
		},
		{
			name: "InternalServerError",
			args: args{
				err: errors.New("another error"),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStatusCode := ErrorStatusCode(tt.args.err); gotStatusCode != tt.wantStatusCode {
				t.Errorf("ErrorStatusCode() = %v, want %v", gotStatusCode, tt.wantStatusCode)
			}
		})
	}
}
