package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"testing"
	"time"
)

var yamlConfig = []byte(`logLevel: debug
logOutput: console

cache:
  ttl: 10
  cleanFrequency: 1
  maxEntries: 600000
  maxEntrySize: 500
  hardMaxCacheSize: 0

invalidator:
  maxWorkers: 5

proxy:
  addr: 0.0.0.0:6081
  backendsAddrs:
    [
      1.2.3.4:5678,
    ]
  response:
    headers:
      set:
        - name: X-Theme
          value: $(resp.header::Theme)

      unset:
        - name: Set-Cookie
          if: $(path) !~ '/preview/' && $(path) !~ '/exit_preview/' && $(cookie::is_preview) == 'True'

  nocache:
    - $(method) == 'POST'
    - $(host) == 'www.kratgo.com'

admin:
  addr: 0.0.0.0:6082
`)

func TestParse(t *testing.T) {
	type args struct {
		filePath    string
		fileContent []byte
		filePerms   os.FileMode
	}

	type want struct {
		err bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Ok",
			args: args{
				filePath:    "/tmp/kratgo_tests.yml",
				fileContent: yamlConfig,
				filePerms:   0775,
			},
			want: want{
				err: false,
			},
		},
		{
			name: "InvalidFile",
			args: args{
				filePath: "",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "InvalidContent",
			args: args{
				filePath:    "/tmp/kratgo_tests.yml",
				fileContent: []byte("aasd\tsxfa\n:$%·$&·&"),
				filePerms:   0775,
			},
			want: want{
				err: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.filePath != "" {
				err := ioutil.WriteFile(tt.args.filePath, tt.args.fileContent, tt.args.filePerms)
				if err != nil {
					panic(err)
				}
			}

			cfg, err := Parse(tt.args.filePath)
			if (err != nil) != tt.want.err {
				t.Fatalf("New() error == '%v', want '%v'", err, tt.want.err)
			}

			if tt.want.err {
				return
			}

			logLevel := "debug"
			if cfg.LogLevel != logLevel {
				t.Fatalf("Parse() LogLevel == '%s', want '%s'", cfg.LogLevel, logLevel)
			}

			logOutput := "console"
			if cfg.LogOutput != logOutput {
				t.Fatalf("Parse() LogOutput == '%s', want '%s'", cfg.LogOutput, logOutput)
			}

			cacheTTL := time.Duration(10)
			if cfg.Cache.TTL != cacheTTL {
				t.Fatalf("Parse() Cache.TTL == '%d', want '%d'", cfg.Cache.TTL, cacheTTL)
			}

			cacheCleanFrequency := time.Duration(1)
			if cfg.Cache.CleanFrequency != cacheCleanFrequency {
				t.Fatalf("Parse() Cache.CleanFrequency == '%d', want '%d'", cfg.Cache.CleanFrequency, cacheCleanFrequency)
			}

			cacheMaxEntries := 600000
			if cfg.Cache.MaxEntries != cacheMaxEntries {
				t.Fatalf("Parse() Cache.MaxEntries == '%d', want '%d'", cfg.Cache.MaxEntries, cacheMaxEntries)
			}

			cacheMaxEntrySize := 500
			if cfg.Cache.MaxEntrySize != cacheMaxEntrySize {
				t.Fatalf("Parse() Cache.MaxEntrySize == '%d', want '%d'", cfg.Cache.MaxEntrySize, cacheMaxEntrySize)
			}

			cacheHardMaxCacheSize := 0
			if cfg.Cache.HardMaxCacheSize != cacheHardMaxCacheSize {
				t.Fatalf("Parse() Cache.HardMaxCacheSize == '%d', want '%d'", cfg.Cache.HardMaxCacheSize, cacheHardMaxCacheSize)
			}

			invalidatorMaxWorkers := int32(5)
			if cfg.Invalidator.MaxWorkers != invalidatorMaxWorkers {
				t.Fatalf("Parse() Invalidator.MaxWorkers == '%d', want '%d'", cfg.Invalidator.MaxWorkers, invalidatorMaxWorkers)
			}

			proxyAddr := "0.0.0.0:6081"
			if cfg.Proxy.Addr != proxyAddr {
				t.Fatalf("Parse() Proxy.Addr == '%s', want '%s'", cfg.Proxy.Addr, proxyAddr)
			}

			proxyBackendsAddrs := []string{"1.2.3.4:5678"}
			if !reflect.DeepEqual(cfg.Proxy.BackendsAddrs, proxyBackendsAddrs) {
				t.Fatalf("Parse() Proxy.BackendsAddrs == '%v', want '%v'", cfg.Proxy.BackendsAddrs, proxyBackendsAddrs)
			}

			proxyResponseHeadersSet := []Header{{Name: "X-Theme", Value: "$(resp.header::Theme)"}}
			if !reflect.DeepEqual(cfg.Proxy.Response.Headers.Set, proxyResponseHeadersSet) {
				t.Fatalf("Parse() Proxy.Response.Headers.Set == '%v', want '%v'", cfg.Proxy.Response.Headers.Set, proxyResponseHeadersSet)
			}

			proxyResponseHeadersUnset := []Header{{Name: "Set-Cookie", When: "$(path) !~ '/preview/' && $(path) !~ '/exit_preview/' && $(cookie::is_preview) == 'True'"}}
			if !reflect.DeepEqual(cfg.Proxy.Response.Headers.Unset, proxyResponseHeadersUnset) {
				t.Fatalf("Parse() Proxy.Response.Headers.Unset == '%v', want '%v'", cfg.Proxy.Response.Headers.Unset, proxyResponseHeadersUnset)
			}

			proxyNocache := []string{"$(method) == 'POST'", "$(host) == 'www.kratgo.com'"}
			if !reflect.DeepEqual(cfg.Proxy.Nocache, proxyNocache) {
				t.Fatalf("Parse() Proxy.Nocache == '%v', want '%v'", cfg.Proxy.Nocache, proxyNocache)
			}

			adminAddr := "0.0.0.0:6082"
			if cfg.Admin.Addr != adminAddr {
				t.Fatalf("Parse() Admin.Addr == '%s', want '%s'", cfg.Admin.Addr, adminAddr)
			}
		})
	}

}

func TestGetEvalParamName(t *testing.T) {
	type args struct {
		key string
	}

	type want struct {
		evalKey      string
		regexEvalKey *regexp.Regexp
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "method",
			args: args{
				key: configMethodVar,
			},
			want: want{
				evalKey: EvalMethodVar,
			},
		},
		{
			name: "host",
			args: args{
				key: configHostVar,
			},
			want: want{
				evalKey: EvalHostVar,
			},
		},
		{
			name: "path",
			args: args{
				key: configPathVar,
			},
			want: want{
				evalKey: EvalPathVar,
			},
		},
		{
			name: "content-type",
			args: args{
				key: configContentTypeVar,
			},
			want: want{
				evalKey: EvalContentTypeVar,
			},
		},
		{
			name: "status-code",
			args: args{
				key: configStatusCodeVar,
			},
			want: want{
				evalKey: EvalStatusCodeVar,
			},
		},
		{
			name: "$(req.header::<NAME>)",
			args: args{
				key: configReqHeaderVar,
			},
			want: want{
				regexEvalKey: regexp.MustCompile(fmt.Sprintf("%s([0-9]{2})", EvalReqHeaderVar)),
			},
		},
		{
			name: "$(resp.header::<NAME>)",
			args: args{
				key: configRespHeaderVar,
			},
			want: want{
				regexEvalKey: regexp.MustCompile(fmt.Sprintf("%s([0-9]{2})", EvalRespHeaderVar)),
			},
		},
		{
			name: "$(cookie::<NAME>)",
			args: args{
				key: configCookieVar,
			},
			want: want{
				regexEvalKey: regexp.MustCompile(fmt.Sprintf("%s([0-9]{2})", EvalCookieVar)),
			},
		},
		{
			name: "unknown",
			args: args{
				key: "$(unknown)",
			},
			want: want{
				evalKey: "$(unknown)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := GetEvalParamName(tt.args.key)

			if tt.want.regexEvalKey != nil {
				if !tt.want.regexEvalKey.MatchString(s) {
					t.Errorf("GetEvalParamName() = '%s', want '%s'", s, tt.want.regexEvalKey.String())
				}
			} else {
				if s != tt.want.evalKey {
					t.Errorf("GetEvalParamName() = '%s', want '%s'", s, tt.want.evalKey)
				}
			}
		})
	}
}

func TestParseConfigKeys(t *testing.T) {
	type args struct {
		key string
	}

	type want struct {
		configKey    string
		evalKey      string
		evalSubKey   string
		regexEvalKey *regexp.Regexp
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "method",
			args: args{
				key: configMethodVar,
			},
			want: want{
				configKey: configMethodVar,
				evalKey:   EvalMethodVar,
			},
		},
		{
			name: "host",
			args: args{
				key: configHostVar,
			},
			want: want{
				configKey: configHostVar,
				evalKey:   EvalHostVar,
			},
		},
		{
			name: "path",
			args: args{
				key: configPathVar,
			},
			want: want{
				configKey: configPathVar,
				evalKey:   EvalPathVar,
			},
		},
		{
			name: "content-type",
			args: args{
				key: configContentTypeVar,
			},
			want: want{
				configKey: configContentTypeVar,
				evalKey:   EvalContentTypeVar,
			},
		},
		{
			name: "status-code",
			args: args{
				key: configStatusCodeVar,
			},
			want: want{
				configKey: configStatusCodeVar,
				evalKey:   EvalStatusCodeVar,
			},
		},
		{
			name: "$(req.header::<NAME>)",
			args: args{
				key: "$(req.header::X-Data)",
			},
			want: want{
				configKey:    "$(req.header::X-Data)",
				evalSubKey:   "X-Data",
				regexEvalKey: regexp.MustCompile(fmt.Sprintf("%s([0-9]{2})", EvalReqHeaderVar)),
			},
		},
		{
			name: "$(resp.header::<NAME>)",
			args: args{
				key: "$(resp.header::X-Data)",
			},
			want: want{
				configKey:    "$(resp.header::X-Data)",
				evalSubKey:   "X-Data",
				regexEvalKey: regexp.MustCompile(fmt.Sprintf("%s([0-9]{2})", EvalRespHeaderVar)),
			},
		},
		{
			name: "$(cookie::<NAME>)",
			args: args{
				key: "$(cookie::Kratgo)",
			},
			want: want{
				configKey:    "$(cookie::Kratgo)",
				evalSubKey:   "Kratgo",
				regexEvalKey: regexp.MustCompile(fmt.Sprintf("%s([0-9]{2})", EvalCookieVar)),
			},
		},
		{
			name: "unknown",
			args: args{
				key: "$(unknown)",
			},
			want: want{
				configKey:  "",
				evalKey:    "",
				evalSubKey: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configKey, evalKey, evalSubKey := ParseConfigKeys(tt.args.key)

			if tt.want.configKey != configKey {
				t.Errorf("GetEvalParamName()[0] = '%s', want '%s'", configKey, tt.want.configKey)
			}

			if tt.want.regexEvalKey != nil {
				if !tt.want.regexEvalKey.MatchString(evalKey) {
					t.Errorf("GetEvalParamName()[1] = '%s', want '%s'", evalKey, tt.want.regexEvalKey.String())
				}
			} else {
				if evalKey != tt.want.evalKey {
					t.Errorf("GetEvalParamName()[1] = '%s', want '%s'", evalKey, tt.want.evalKey)
				}
			}

			if tt.want.evalSubKey != evalSubKey {
				t.Errorf("GetEvalParamName()[2] = '%s', want '%s'", evalSubKey, tt.want.evalSubKey)
			}
		})
	}
}
