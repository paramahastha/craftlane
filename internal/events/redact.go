package events

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
)

var secretKey = regexp.MustCompile(`(?i)(token|secret|password|passwd|api[_-]?key|authorization|private[_-]?key|credential)`)
var secretValue = regexp.MustCompile(`(?i)(bearer\s+|ghp_|github_pat_|sk-[A-Za-z0-9_-]+|AKIA[0-9A-Z]{16})[A-Za-z0-9_./+=-]*`)

func Redact(v any) any {
	switch x := v.(type) {
	case map[string]any:
		o := map[string]any{}
		for k, v := range x {
			if secretKey.MatchString(k) {
				o[k] = "[REDACTED]"
			} else {
				o[k] = Redact(v)
			}
		}
		return o
	case []any:
		o := make([]any, len(x))
		for i, v := range x {
			o[i] = Redact(v)
		}
		return o
	case string:
		return secretValue.ReplaceAllString(x, "[REDACTED]")
	default:
		kind := reflect.TypeOf(v)
		if kind != nil && (kind.Kind() == reflect.Map || kind.Kind() == reflect.Slice || kind.Kind() == reflect.Array || kind.Kind() == reflect.Struct || kind.Kind() == reflect.Ptr) {
			if data, err := json.Marshal(v); err == nil {
				var plain any
				if json.Unmarshal(data, &plain) == nil {
					return Redact(plain)
				}
			}
		}
		return v
	}
}
func RedactText(s string) string {
	return strings.TrimSpace(secretValue.ReplaceAllString(s, "[REDACTED]"))
}
