package admin

import (
	"testing"

	"github.com/savsgio/kratgo/modules/invalidator"

	"github.com/savsgio/atreugo/v11"
	"github.com/valyala/fasthttp"
)

func TestAdmin_invalidateView(t *testing.T) {
	type args struct {
		method   string
		body     string
		addError error
	}

	type want struct {
		response   string
		statusCode int
		err        bool
		callAdd    bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Ok",
			args: args{
				method: "POST",
				body:   "{\"host\": \"www.kratgo.com\"}",
			},
			want: want{
				response:   "OK",
				statusCode: 200,
				err:        false,
				callAdd:    true,
			},
		},
		{
			name: "EmptyJSONBody",
			args: args{
				method:   "POST",
				body:     "{}",
				addError: invalidator.ErrEmptyFields,
			},
			want: want{
				response:   invalidator.ErrEmptyFields.Error(),
				statusCode: 400,
				err:        false,
				callAdd:    false,
			},
		},
		{
			name: "InvalidJSONBody",
			args: args{
				method: "POST",
				body:   "\"",
			},
			want: want{
				response:   "",  // err message is setted by Atreugo when is returned the error
				statusCode: 200, // 500 is setted by Atreugo when is returned the error
				err:        true,
				callAdd:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invalidatorMock := &mockInvalidator{
				err: tt.args.addError,
			}

			admin, err := New(testConfig())
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			admin.invalidator = invalidatorMock

			actx := new(atreugo.RequestCtx)
			actx.RequestCtx = new(fasthttp.RequestCtx)

			actx.Request.Header.SetMethod(tt.args.method)
			actx.Request.SetBodyString(tt.args.body)

			err = admin.invalidateView(actx)
			if (err != nil) != tt.want.err {
				t.Fatalf("Admin.invalidateView() error == '%v', want '%v'", err, tt.want.err)
			}

			if tt.want.callAdd && !invalidatorMock.addCalled {
				t.Error("Admin.invalidateView() has not called to admin.invalidator.Add(...)")
			}

			statusCode := actx.Response.StatusCode()
			if statusCode != tt.want.statusCode {
				t.Errorf("Admin.invalidateView() status code == '%d', want '%d'", statusCode, tt.want.statusCode)
			}

			respBody := string(actx.Response.Body())
			if respBody != tt.want.response {
				t.Errorf("Admin.invalidateView() response body == '%s', want '%s'", respBody, tt.want.response)
			}
		})
	}
}
