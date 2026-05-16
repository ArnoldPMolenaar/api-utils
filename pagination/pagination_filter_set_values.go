package pagination

import (
	"slices"
	"strings"

	"github.com/valyala/fasthttp"
)

// FilterSetValues adds or updates values for a specific column in the pagination filter arguments.
func FilterSetValues(values *fasthttp.Args, param, column string, valuesToAdd ...string) {
	if values == nil || len(valuesToAdd) == 0 || param == "" || column == "" {
		return
	}
	if !slices.Contains([]string{"searchLike", "searchEq", "searchLikeOr", "searchEqOr", "searchIn", "searchBetween"}, param) {
		return
	}

	value := column + ":" + strings.Join(valuesToAdd, ";")
	existingValue := values.Peek(param)

	if existingValue != nil {
		values.Set(param, string(existingValue)+","+value)
	} else {
		values.Set(param, value)
	}
}
