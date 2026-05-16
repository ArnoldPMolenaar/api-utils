package pagination

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func TestFilterSetValues(t *testing.T) {
	t.Run("no-op on invalid input", func(t *testing.T) {
		FilterSetValues(nil, "searchEq", "tenant_id", "1")

		args := &fasthttp.Args{}
		FilterSetValues(args, "", "tenant_id", "1")
		FilterSetValues(args, "searchEq", "", "1")
		FilterSetValues(args, "searchEq", "tenant_id")
		FilterSetValues(args, "unknown", "tenant_id", "1")

		if got := string(args.Peek("searchEq")); got != "" {
			t.Fatalf("expected searchEq to remain empty, got %q", got)
		}
	})

	t.Run("sets new value", func(t *testing.T) {
		args := &fasthttp.Args{}
		FilterSetValues(args, "searchEq", "tenant_id", "42")

		if got := string(args.Peek("searchEq")); got != "tenant_id:42" {
			t.Fatalf("expected tenant_id:42, got %q", got)
		}
	})

	t.Run("appends to existing value", func(t *testing.T) {
		args := &fasthttp.Args{}
		args.Set("searchEq", "status:active")

		FilterSetValues(args, "searchEq", "tenant_id", "1", "2")

		if got := string(args.Peek("searchEq")); got != "status:active,tenant_id:1;2" {
			t.Fatalf("expected appended value, got %q", got)
		}
	})
}

