package pagination

import (
	"reflect"
	"strings"
	"testing"

	"gorm.io/gorm"
)

func newTestDB() *gorm.DB {
	return &gorm.DB{Config: &gorm.Config{}}
}

func assertDBErrorContains(t *testing.T, db *gorm.DB, want string) {
	t.Helper()
	if db.Error == nil {
		t.Fatalf("expected db error containing %q, got nil", want)
	}
	if !strings.Contains(db.Error.Error(), want) {
		t.Fatalf("expected db error containing %q, got %q", want, db.Error.Error())
	}
}

func TestParseColumn(t *testing.T) {
	cases := []struct {
		in   string
		out  string
		desc string
	}{
		{in: "column", out: "\"column\"", desc: "simple column"},
		{in: "table.column", out: "\"table\".\"column\"", desc: "table column"},
		{in: "schema.table.column", out: "\"schema\".\"table\".\"column\"", desc: "schema table column"},
		{in: "  table .  column  ", out: "\"table\".\"column\"", desc: "spaces around"},
		{in: "\"table\".\"column\"", out: "\"table\".\"column\"", desc: "already quoted"},
		{in: "", out: "\"\"", desc: "empty string"},
		{in: ".column", out: "\"\".\"column\"", desc: "leading dot"},
		{in: "table.", out: "\"table\".\"\"", desc: "trailing dot"},
		{in: "table..column", out: "\"table\".\"\".\"column\"", desc: "double dot"},
	}

	for _, c := range cases {
		got := parseColumn(c.in)
		if got != c.out {
			t.Fatalf("%s: parseColumn(%q) = %q, want %q", c.desc, c.in, got, c.out)
		}
	}
}

func TestCount(t *testing.T) {
	if got := Count(101, 10); got != 11 {
		t.Fatalf("Count(101, 10) = %d, want 11", got)
	}
	if got := Count(100, 10); got != 10 {
		t.Fatalf("Count(100, 10) = %d, want 10", got)
	}
}

func TestOffset(t *testing.T) {
	if got := Offset(1, 25); got != 0 {
		t.Fatalf("Offset(1, 25) = %d, want 0", got)
	}
	if got := Offset(3, 25); got != 50 {
		t.Fatalf("Offset(3, 25) = %d, want 50", got)
	}
}

func TestCreatePaginationModel(t *testing.T) {
	rows := []string{"a", "b"}
	got := CreatePaginationModel(10, 2, 5, 42, rows)

	if got.Limit != 10 || got.Page != 2 || got.PageCount != 5 || got.Total != 42 {
		t.Fatalf("unexpected metadata in model: %+v", got)
	}

	res, ok := got.Result.([]string)
	if !ok {
		t.Fatalf("expected []string result, got %T", got.Result)
	}
	if !reflect.DeepEqual(res, rows) {
		t.Fatalf("unexpected result: got %v, want %v", res, rows)
	}
}

func TestParseSingleValueParams(t *testing.T) {
	allowed := map[string]bool{"name": true, "status": true}

	t.Run("parses valid pairs", func(t *testing.T) {
		db := newTestDB()
		got := parseSingleValueParams(db, "name:john,status:active", allowed)
		want := map[string]string{"name": "john", "status": "active"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("parseSingleValueParams() = %v, want %v", got, want)
		}
		if db.Error != nil {
			t.Fatalf("expected nil db error, got %v", db.Error)
		}
	})

	t.Run("ignores empty value", func(t *testing.T) {
		db := newTestDB()
		got := parseSingleValueParams(db, "name:", allowed)
		if len(got) != 0 {
			t.Fatalf("expected empty map, got %v", got)
		}
		if db.Error != nil {
			t.Fatalf("expected nil db error, got %v", db.Error)
		}
	})

	t.Run("invalid format adds error", func(t *testing.T) {
		db := newTestDB()
		_ = parseSingleValueParams(db, "name", allowed)
		assertDBErrorContains(t, db, "cannot parse invalid format")
	})

	t.Run("disallowed column adds error", func(t *testing.T) {
		db := newTestDB()
		got := parseSingleValueParams(db, "role:admin", allowed)
		if len(got) != 0 {
			t.Fatalf("expected empty map, got %v", got)
		}
		assertDBErrorContains(t, db, "column not allowed")
	})
}

func TestParseMultiValueParams(t *testing.T) {
	allowed := map[string]bool{"id": true, "status": true}

	t.Run("parses valid pairs", func(t *testing.T) {
		db := newTestDB()
		got := parseMultiValueParams(db, "id:1;2;3,status:active;disabled", allowed)
		want := map[string][]string{
			"id":     []string{"1", "2", "3"},
			"status": []string{"active", "disabled"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("parseMultiValueParams() = %v, want %v", got, want)
		}
		if db.Error != nil {
			t.Fatalf("expected nil db error, got %v", db.Error)
		}
	})

	t.Run("ignores empty value", func(t *testing.T) {
		db := newTestDB()
		got := parseMultiValueParams(db, "id:", allowed)
		if len(got) != 0 {
			t.Fatalf("expected empty map, got %v", got)
		}
		if db.Error != nil {
			t.Fatalf("expected nil db error, got %v", db.Error)
		}
	})

	t.Run("invalid format adds error", func(t *testing.T) {
		db := newTestDB()
		_ = parseMultiValueParams(db, "id", allowed)
		assertDBErrorContains(t, db, "cannot parse invalid format")
	})

	t.Run("disallowed column adds error", func(t *testing.T) {
		db := newTestDB()
		got := parseMultiValueParams(db, "role:admin;user", allowed)
		if len(got) != 0 {
			t.Fatalf("expected empty map, got %v", got)
		}
		assertDBErrorContains(t, db, "column not allowed")
	})
}

