package simpleapi

import (
	"errors"
	"github.com/redagain/go-testrequest"
	"net/http"
	"reflect"
	"testing"
)

func Test_bearerAuthMiddleware_HandleRequest(t *testing.T) {
	type fields struct {
		introspectToken func(accessToken string) (active bool, err error)
	}
	type args struct {
		w    http.ResponseWriter
		req  *http.Request
		next RequestHandler
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name:   "InvalidAuth",
			fields: fields{},
			args: args{
				w:   testrequest.NopResponseWriter(),
				req: testrequest.Builder().Request(),
			},
			wantErr: ErrUnauthorized,
		},
		{
			name: "TokenIntrospectionError",
			fields: fields{
				introspectToken: func(accessToken string) (active bool, err error) {
					err = errors.New("failed to verify the token")
					return
				},
			},
			args: args{
				w:   testrequest.NopResponseWriter(),
				req: testrequest.Builder().SetBearerAuth("CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ").Request(),
			},
			wantErr: errors.New("failed to verify the token"),
		},
		{
			name: "TokenNotActive",
			fields: fields{
				introspectToken: func(accessToken string) (active bool, err error) {
					return
				},
			},
			args: args{
				w:   testrequest.NopResponseWriter(),
				req: testrequest.Builder().SetBearerAuth("CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ").Request(),
			},
			wantErr: ErrUnauthorized,
		},
		{
			name: "NextHandler",
			fields: fields{
				introspectToken: func(accessToken string) (active bool, err error) {
					active = true
					return
				},
			},
			args: args{
				w:   testrequest.NopResponseWriter(),
				req: testrequest.Builder().SetBearerAuth("CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ").Request(),
				next: func(w http.ResponseWriter, req *http.Request) (err error) {
					return
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &bearerAuthMiddleware{
				introspectToken: tt.fields.introspectToken,
			}
			err := h.HandleRequest(tt.args.w, tt.args.req, tt.args.next)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("HandleRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
