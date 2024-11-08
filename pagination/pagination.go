package pagination

import (
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
)

// Model struct is used to return paginated data.
type Model struct {
	Limit     int         `json:"limit"`
	Page      int         `json:"page"`
	PageCount int         `json:"pageCount"`
	Total     int         `json:"total"`
	Result    interface{} `json:"result"`
}

// Query builds a pagination query with the provided values
// and checks the input columns against the allowedColumns list.
// Returns a gorm query to be used in the function or an error.
func Query(args *fasthttp.Args, allowedColumns map[string]bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = parseSearchLike(args.Peek("searchLike"), db, allowedColumns)
		db = parseSearchEq(args.Peek("searchEq"), db, allowedColumns)
		db = parseSearchLikeOr(args.Peek("searchLikeOr"), db, allowedColumns)
		db = parseSearchEqOr(args.Peek("searchEqOr"), db, allowedColumns)
		db = parseSearchIn(args.Peek("searchIn"), db, allowedColumns)
		db = parseSearchBetween(args.Peek("searchBetween"), db, allowedColumns)

		return db
	}
}

// Sort builds a sort query with the provided values
// and checks the input columns against the allowedColumns list.
// Returns a gorm query to be used in the function or an error.
func Sort(args *fasthttp.Args, allowedColumns map[string]bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = parseSortBy(args.Peek("sortBy"), db, allowedColumns)

		return db
	}
}

// Count calculates the page count with the given resultCount of a pagination query and a page limit.
func Count(resultCount, limit int) int {
	return int(math.Ceil(float64(resultCount) / float64(limit)))
}

// Offset calculates the offset with the page and limit params
func Offset(page, limit int) int {
	return (page - 1) * limit
}

// CreatePaginationModel is a helper to be able to return a pagination model in a single line
func CreatePaginationModel(limit, page, pageCount, total int, result interface{}) Model {
	return Model{
		Limit:     limit,
		Page:      page,
		PageCount: pageCount,
		Total:     total,
		Result:    result,
	}
}

// search_like: for |where ... LIKE ... AND| query = search_like=column:value,column:value =>
// search_like=firstname:john,lastname:doe
func parseSearchLike(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		db = db.Where(fmt.Sprintf("%s ILIKE ?", key), fmt.Sprintf("%%%s%%", value))
	}

	return db
}

// search_eq: for |where ... = ... AND| query = search_eq=column:value,column:value =>
// search_eq=firstname:john,lastname:doe
func parseSearchEq(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}

	return db
}

// search_like_or: for |where ... like ... OR| query = search_like_or=column:value,column:value =>
// search_or_like=firstname:john,lastname:doe
func parseSearchLikeOr(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	var conditions []string
	var values []interface{}
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", key))
		values = append(values, fmt.Sprintf("%%%s%%", value))
	}

	if len(conditions) > 0 {
		db = db.Where(strings.Join(conditions, " OR "), values...)
	}

	return db
}

// search_eq_or: for |where ... = ... OR| query = search_eq_or=column:value,column:value =>
// search_or_eq=firstname:john,lastname:doe
func parseSearchEqOr(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	var conditions []string
	var values []interface{}
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		conditions = append(conditions, fmt.Sprintf("%s = ?", key))
		values = append(values, value)
	}

	if len(conditions) > 0 {
		db = db.Where(strings.Join(conditions, " OR "), values...)
	}

	return db
}

// search_in: for |where IN| query = search_in=column:value.value.value => search_in=is_online:true.false
func parseSearchIn(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseMultiValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		db = db.Where(fmt.Sprintf("%s IN (?)", key), value)
	}

	return db
}

// search_between  for |where ... between ... AND ...| query = search_between=column:value1.value2 =>
// search_between=created_at:2020-08-03.2020-09-03
func parseSearchBetween(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseMultiValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		if len(value) != 2 {
			_ = db.AddError(errors.New("not exactly two values for between query"))
		}

		// Parse the date-time strings
		startTime, err1 := time.Parse(time.RFC3339, value[0])
		endTime, err2 := time.Parse(time.RFC3339, value[1])
		if err1 != nil || err2 != nil {
			_ = db.AddError(errors.New("invalid date-time format"))
		}

		db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", key), startTime, endTime)
	}

	return db
}

// sort_by: for |ORDER BY| query = sort_by=column:value,column:value => sort_by=firstname:asc,lastname:desc
func parseSortBy(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		switch value {
		case "desc":
			db = db.Order(fmt.Sprintf("%s DESC", key))
		case "asc":
			db = db.Order(fmt.Sprintf("%s ASC", key))
		default:
			_ = db.AddError(errors.New("order not asc or desc"))
		}
	}

	return db
}

// parseSingleValueParams parses the query string for single value params.
// The query string should be in the format of key:value,key:value
func parseSingleValueParams(db *gorm.DB, params string, allowedColumns map[string]bool) map[string]string {
	paramMap := make(map[string]string)

	if params != "" {
		paramSearchParts := strings.Split(params, ",")
		for _, paramSearchPart := range paramSearchParts {
			valuePairs := strings.Split(paramSearchPart, ":")
			canParse := len(valuePairs) == 2 && valuePairs[0] != "" && valuePairs[1] != ""
			isAllowed := allowedColumns[valuePairs[0]]

			if !canParse {
				_ = db.AddError(errors.New("cannot parse invalid format"))
			}
			if !isAllowed {
				_ = db.AddError(errors.New("column not allowed"))
			}
			if isAllowed && canParse {
				paramMap[valuePairs[0]] = valuePairs[1]
			}
		}
	}

	return paramMap
}

// parseMultiValueParams parses the query string for multi value params.
// The query string should be in the format of key:value.value.value,key:value.value.value
func parseMultiValueParams(db *gorm.DB, params string, allowedColumns map[string]bool) map[string][]string {
	paramMap := make(map[string][]string)

	if params != "" {
		paramSearchParts := strings.Split(params, ",")
		for _, paramSearchPart := range paramSearchParts {
			valuePairs := strings.SplitN(paramSearchPart, ":", 2)
			canParse := len(valuePairs) == 2 && valuePairs[0] != "" && valuePairs[1] != ""
			isAllowed := allowedColumns[valuePairs[0]]

			if !canParse {
				_ = db.AddError(errors.New("cannot parse invalid format"))
			}
			if !isAllowed {
				_ = db.AddError(errors.New("column not allowed"))
			}
			if isAllowed && canParse {
				paramMap[valuePairs[0]] = strings.Split(valuePairs[1], ";")
			}
		}
	}

	return paramMap
}
