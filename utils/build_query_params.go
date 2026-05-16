package utils

import "fmt"

// BuildQueryParams builds a query string with repeated keys for each value.
// Example: BuildQueryParams("appName", []string{"foo", "bar"}) returns "appName=foo&appName=bar"
func BuildQueryParams(key string, values []string) string {
	var params string

	for i, v := range values {
		if i > 0 {
			params += "&"
		}
		params += fmt.Sprintf("%s=%s", key, v)
	}

	return params
}
