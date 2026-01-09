package pagination

import "testing"

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
