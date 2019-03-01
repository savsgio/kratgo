package kratgo

import (
	"testing"
)

func Test_getLogOutput(t *testing.T) {
	type args struct {
		output string
	}

	type want struct {
		fileName string
		err      bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "console",
			args: args{output: defaultLogOutput},
			want: want{
				fileName: "/dev/stderr",
				err:      false,
			},
		},
		{
			name: "file",
			args: args{output: "/tmp/kratgo_test.log"},
			want: want{
				fileName: "/tmp/kratgo_test.log",
				err:      false,
			},
		},
		{
			name: "invalid",
			args: args{output: ""},
			want: want{
				err: true,
			},
		},
		{
			name: "error",
			args: args{output: "/sadasdadr2343dcr4c234/kratgo_test.log"},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := getLogOutput(tt.args.output)
			if (err != nil) != tt.want.err {
				t.Fatalf("getLogOutput() unexpected error: %v", err)
			}

			if f != nil {
				fileName := f.Name()
				if fileName != tt.want.fileName {
					t.Errorf("getLogOutput() file = '%s', want '%s'", fileName, tt.want.fileName)
				}
			}
		})
	}
}
