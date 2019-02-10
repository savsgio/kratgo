package config

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"

	yaml "gopkg.in/yaml.v2"
)

var configEvaluationVars = map[string]string{
	configMethodVar:      EvalMethodVar,
	configHostVar:        EvalHostVar,
	configPathVar:        EvalPathVar,
	configContentTypeVar: EvalContentTypeVar,
	configStatusCodeVar:  EvalStatusCodeVar,
	configReqHeaderVar:   EvalReqHeaderVar,
	configRespHeaderVar:  EvalRespHeaderVar,
	configCookieVar:      EvalCookieVar,
}

// ConfigVarRegex ...
var ConfigVarRegex = regexp.MustCompile("\\$\\([a-zA-Z0-9\\-:\\.\\_]+\\)")

// ConfigReqHeaderVarRegex ...
var ConfigReqHeaderVarRegex = regexp.MustCompile("\\$\\(req\\.header::([a-zA-Z0-9\\-\\_]+)\\)")

// ConfigRespHeaderVarRegex ...
var ConfigRespHeaderVarRegex = regexp.MustCompile("\\$\\(resp\\.header::([a-zA-Z0-9\\-\\_]+)\\)")

// ConfigCookieVarRegex ...
var ConfigCookieVarRegex = regexp.MustCompile("\\$\\(cookie::([a-zA-Z0-9\\-\\_]+)\\)")

// Parse ...
func Parse(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetEvalParamName ...
func GetEvalParamName(k string) string {
	if k == configReqHeaderVar || k == configRespHeaderVar {
		return fmt.Sprintf("%s%d", configEvaluationVars[k], rand.Int31n(100))

	} else if k == configCookieVar {
		return fmt.Sprintf("%s%d", configEvaluationVars[k], rand.Int31n(100))
	}

	if v, ok := configEvaluationVars[k]; ok {
		return v
	}

	return k
}

// ParseConfigKeys ...
func ParseConfigKeys(s string) (string, string, string) {
	for k := range configEvaluationVars {
		if k == configReqHeaderVar {
			data := ConfigReqHeaderVarRegex.FindStringSubmatch(s)
			if len(data) > 1 {
				return data[0], GetEvalParamName(k), data[1]
			}
		} else if k == configRespHeaderVar {
			data := ConfigRespHeaderVarRegex.FindStringSubmatch(s)
			if len(data) > 1 {
				return data[0], GetEvalParamName(k), data[1]
			}

		} else if k == configCookieVar {
			data := ConfigCookieVarRegex.FindStringSubmatch(s)
			if len(data) > 1 {
				return data[0], GetEvalParamName(k), data[1]
			}

		} else {
			data := ConfigVarRegex.FindStringSubmatch(s)
			if len(data) > 0 && data[0] == k {
				return data[0], GetEvalParamName(data[0]), ""
			}
		}
	}

	return "", "", ""
}
