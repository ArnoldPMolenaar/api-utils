package pagination

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
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
		// Combine the OR groups (EqOr and LikeOr) into a single OR clause that is AND-ed with other filters
		db = parseOr(args.Peek("searchEqOr"), args.Peek("searchLikeOr"), db, allowedColumns)
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

// parseSearchLike Adds LIKE conditions to the GORM DB query
// searchLike: for |where ... LIKE ... AND| query = searchLike=column:value,column:value =>
// searchLike=firstname:john,lastname:doe
func parseSearchLike(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		db = db.Where(fmt.Sprintf("CAST(%s AS TEXT) ILIKE ?", parseColumn(key)), fmt.Sprintf("%%%s%%", value))
	}

	return db
}

// parseSearchEq Adds equality conditions to the GORM DB query
// searchEq: for |where ... = ... AND| query = searchEq=column:value,column:value =>
// searchEq=firstname:john,lastname:doe
func parseSearchEq(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		db = db.Where(fmt.Sprintf("CAST(%s AS TEXT) = ?", parseColumn(key)), value)
	}

	return db
}

// parseOr merges searchEqOr and searchLikeOr into a single OR group that is AND-ed with other filters.
// Example: searchEqOr=a:1,b:2 and searchLikeOr=c:x => WHERE (... AND (... OR ... OR ...))
func parseOr(eqParams []byte, likeParams []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	var conditions []string
	var values []interface{}

	// Equal OR part
	eqMap := parseSingleValueParams(db, string(eqParams), allowedColumns)
	for key, value := range eqMap {
		conditions = append(conditions, fmt.Sprintf("CAST(%s AS TEXT) = ?", parseColumn(key)))
		values = append(values, value)
	}

	// LIKE OR part
	likeMap := parseSingleValueParams(db, string(likeParams), allowedColumns)
	for key, value := range likeMap {
		conditions = append(conditions, fmt.Sprintf("CAST(%s AS TEXT) ILIKE ?", parseColumn(key)))
		values = append(values, fmt.Sprintf("%%%s%%", value))
	}

	if len(conditions) > 0 {
		group := "(" + strings.Join(conditions, " OR ") + ")"
		db = db.Where(group, values...)
	}

	return db
}

// parseSearchIn Adds IN conditions to the GORM DB query
// searchIn: for |where IN| query = searchIn=column:value;value;value => searchIn=is_online:true;false
func parseSearchIn(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseMultiValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		db = db.Where(fmt.Sprintf("CAST(%s AS TEXT) IN (?)", parseColumn(key)), value)
	}

	return db
}

// parseSearchBetween Adds BETWEEN conditions to the GORM DB query
// searchBetween: for |where ... between ... AND ...| query = searchBetween=column:value1;value2 =>
// searchBetween=created_at:2020-08-03;2020-09-03
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

		db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", parseColumn(key)), startTime, endTime)
	}

	return db
}

// parseSortBy Adds ORDER BY conditions to the GORM DB query
// sortBy: for |ORDER BY| query = sortBy=column:value,column:value => sortBy=firstname:asc,lastname:desc
func parseSortBy(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		switch value {
		case "desc":
			db = db.Order(fmt.Sprintf("%s DESC", parseColumn(key)))
		case "asc":
			db = db.Order(fmt.Sprintf("%s ASC", parseColumn(key)))
		default:
			_ = db.AddError(errors.New("order not asc or desc"))
		}
	}

	return db
}

// parseColumn quotes SQL identifiers correctly for GORM/SQL.
// It splits the input on dots and wraps each non-empty part with double quotes,
// so:
//   - "table.column" => "table"."column"
//   - "schema.table.column" => "schema"."table"."column"
//
// The function trims whitespace and strips any existing double quotes first to
// avoid double-quoting, then re-quotes consistently. Empty parts (e.g., due to
// leading/trailing/double dots) are converted to "" to preserve structure,
// allowing upstream validation to fail fast if inputs are malformed.
func parseColumn(column string) string {
	// Normalize: trim spaces and remove existing quotes to avoid double-quoting or malformed input
	cleaned := strings.TrimSpace(column)
	if cleaned == "" {
		return "\"\""
	}
	cleaned = strings.ReplaceAll(cleaned, "\"", "")

	parts := strings.Split(cleaned, ".")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			// Preserve position to keep dot structure, although it will likely lead to a SQL error upstream
			parts[i] = "\"\""
			continue
		}
		parts[i] = fmt.Sprintf("\"%s\"", p)
	}
	return strings.Join(parts, ".")
}

// parseSingleValueParams parses the query string for single value params.
// The query string should be in the format of key:value,key:value
func parseSingleValueParams(db *gorm.DB, params string, allowedColumns map[string]bool) map[string]string {
	paramMap := make(map[string]string)

	if params != "" {
		paramSearchParts := strings.Split(params, ",")
		for _, paramSearchPart := range paramSearchParts {
			valuePairs := strings.Split(paramSearchPart, ":")
			// malformed when not exactly 2 parts or key empty
			malformed := len(valuePairs) != 2 || valuePairs[0] == ""
			if malformed {
				_ = db.AddError(errors.New("cannot parse invalid format"))
				continue
			}

			key := valuePairs[0]
			val := valuePairs[1]

			// skip silently if value is empty (no error, just ignore this pair)
			if val == "" {
				continue
			}

			isAllowed := allowedColumns[key]
			if !isAllowed {
				_ = db.AddError(errors.New("column not allowed"))
				continue
			}

			paramMap[key] = val
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
			// malformed when not exactly 2 parts or key empty
			malformed := len(valuePairs) != 2 || valuePairs[0] == ""
			if malformed {
				_ = db.AddError(errors.New("cannot parse invalid format"))
				continue
			}

			key := valuePairs[0]
			val := valuePairs[1]

			// skip silently if value part is empty (no values)
			if val == "" {
				continue
			}

			isAllowed := allowedColumns[key]
			if !isAllowed {
				_ = db.AddError(errors.New("column not allowed"))
				continue
			}

			paramMap[key] = strings.Split(val, ";")
		}
	}

	return paramMap
}
