package pagination

import (
	"reflect"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestFilterExtractValues(t *testing.T) {
	t.Run("nil args", func(t *testing.T) {
		values, param, query := FilterExtractValues(nil, "tenant_id")
		if len(values) != 0 {
			t.Fatalf("expected no values, got %v", values)
		}
		if param != "" {
			t.Fatalf("expected empty param, got %q", param)
		}
		if query != nil {
			t.Fatalf("expected nil query, got %#v", query)
		}
	})

	t.Run("extracts and removes single value", func(t *testing.T) {
		args := &fasthttp.Args{}
		args.Set("searchEq", "tenant_id:42,status:active")
		args.Set("page", "2")

		values, param, query := FilterExtractValues(args, "tenant_id")
		if !reflect.DeepEqual(values, []string{"42"}) {
			t.Fatalf("expected [42], got %v", values)
		}
		if param != "searchEq" {
			t.Fatalf("expected param searchEq, got %q", param)
		}
		if got := string(query.Peek("searchEq")); got != "status:active" {
			t.Fatalf("expected remaining searchEq to be status:active, got %q", got)
		}
		if got := string(query.Peek("page")); got != "2" {
			t.Fatalf("expected unrelated param to remain, got %q", got)
		}
	})

	t.Run("extracts semicolon separated values and deletes empty filter", func(t *testing.T) {
		args := &fasthttp.Args{}
		args.Set("searchIn", "tenant_id:1;2;3")

		values, param, query := FilterExtractValues(args, "tenant_id")
		if !reflect.DeepEqual(values, []string{"1", "2", "3"}) {
			t.Fatalf("expected [1 2 3], got %v", values)
		}
		if param != "searchIn" {
			t.Fatalf("expected param searchIn, got %q", param)
		}
		if got := string(query.Peek("searchIn")); got != "" {
			t.Fatalf("expected searchIn to be removed, got %q", got)
		}
	})

	t.Run("extracts and removes trailing value without leaving comma", func(t *testing.T) {
		args := &fasthttp.Args{}
		args.Set("searchEq", "status:active,tenant_id:42")

		values, param, query := FilterExtractValues(args, "tenant_id")
		if !reflect.DeepEqual(values, []string{"42"}) {
			t.Fatalf("expected [42], got %v", values)
		}
		if param != "searchEq" {
			t.Fatalf("expected param searchEq, got %q", param)
		}
		if got := string(query.Peek("searchEq")); got != "status:active" {
			t.Fatalf("expected remaining searchEq to be status:active, got %q", got)
		}
	})

	t.Run("ignores unsupported params", func(t *testing.T) {
		args := &fasthttp.Args{}
		args.Set("custom", "tenant_id:42")

		values, param, query := FilterExtractValues(args, "tenant_id")
		if len(values) != 0 {
			t.Fatalf("expected no extracted values, got %v", values)
		}
		if param != "" {
			t.Fatalf("expected empty param, got %q", param)
		}
		if got := string(query.Peek("custom")); got != "tenant_id:42" {
			t.Fatalf("expected custom param unchanged, got %q", got)
		}
	})
}


