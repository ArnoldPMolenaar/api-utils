package pagination

import (
	"slices"
	"strings"

	"github.com/valyala/fasthttp"
)

// FilterExtractValues extracts the values for a specific column from the pagination filter arguments.
//
//	values: needs to be a querystring in the format "?column_name:value1,value2,...;other_column:other_value&"
func FilterExtractValues(values *fasthttp.Args, column string) (result []string, param string, query *fasthttp.Args) {
	if values == nil {
		return result, param, values
	}

	newValues := &fasthttp.Args{}
	values.CopyTo(newValues)

	for key, value := range values.All() {
		k := string(key)
		val := string(value)
		if val == "" || !slices.Contains([]string{"searchLike", "searchEq", "searchLikeOr", "searchEqOr", "searchIn", "searchBetween"}, k) {
			continue
		}

		paramSearchParts := strings.Split(val, ",")
		for _, paramSearchPart := range paramSearchParts {
			valuePairs := strings.Split(paramSearchPart, ":")
			if valuePairs[0] != column {
				continue
			}

			if strings.Contains(valuePairs[1], ";") {
				result = append(result, strings.Split(valuePairs[1], ";")...)
			} else {
				result = append(result, valuePairs[1])
			}

			param = k
			newVal := strings.Replace(val, paramSearchPart, "", 1)
			newVal = strings.Replace(newVal, ",,", ",", -1)
			newVal = strings.Trim(newVal, ",")

			if newVal == "" {
				newValues.Del(k)
			} else {
				newValues.Set(k, newVal)
			}

			break
		}
	}

	return result, param, newValues
}
