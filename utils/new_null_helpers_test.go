package utils

import (
	"database/sql"
	"testing"
	"time"
)

func TestNewNullGenericAndPtrFromNull(t *testing.T) {
	v := 7
	n := NewNull(&v)
	if !n.Valid || n.V != 7 {
		t.Fatalf("unexpected sql.Null value: %+v", n)
	}
	ptr := PtrFromNull(n)
	if ptr == nil || *ptr != 7 {
		t.Fatalf("unexpected pointer from null: %v", ptr)
	}

	nNil := NewNull[int](nil)
	if nNil.Valid {
		t.Fatalf("expected invalid null from nil pointer: %+v", nNil)
	}
	if PtrFromNull(nNil) != nil {
		t.Fatal("expected nil pointer from invalid null")
	}
}

func TestSpecificNullHelpers(t *testing.T) {
	b := true
	nb := NewNullBool(&b)
	if !nb.Valid || nb.Bool != b || PtrFromNullBool(nb) == nil || *PtrFromNullBool(nb) != b {
		t.Fatalf("unexpected null bool behavior: %+v", nb)
	}
	if NewNullBool(nil).Valid || PtrFromNullBool(sql.NullBool{}) != nil {
		t.Fatal("expected invalid/null behavior for null bool")
	}

	by := byte(8)
	nby := NewNullByte(&by)
	if !nby.Valid || nby.Byte != by || PtrFromNullByte(nby) == nil || *PtrFromNullByte(nby) != by {
		t.Fatalf("unexpected null byte behavior: %+v", nby)
	}

	f := 3.14
	nf := NewNullFloat64(&f)
	if !nf.Valid || nf.Float64 != f || PtrFromNullFloat64(nf) == nil || *PtrFromNullFloat64(nf) != f {
		t.Fatalf("unexpected null float behavior: %+v", nf)
	}

	i16 := int16(16)
	n16 := NewNullInt16(&i16)
	if !n16.Valid || n16.Int16 != i16 || PtrFromNullInt16(n16) == nil || *PtrFromNullInt16(n16) != i16 {
		t.Fatalf("unexpected null int16 behavior: %+v", n16)
	}

	i32 := int32(32)
	n32 := NewNullInt32(&i32)
	if !n32.Valid || n32.Int32 != i32 || PtrFromNullInt32(n32) == nil || *PtrFromNullInt32(n32) != i32 {
		t.Fatalf("unexpected null int32 behavior: %+v", n32)
	}

	i64 := int64(64)
	n64 := NewNullInt64(&i64)
	if !n64.Valid || n64.Int64 != i64 || PtrFromNullInt64(n64) == nil || *PtrFromNullInt64(n64) != i64 {
		t.Fatalf("unexpected null int64 behavior: %+v", n64)
	}

	s := "hello"
	ns := NewNullString(&s)
	if !ns.Valid || ns.String != s || PtrFromNullString(ns) == nil || *PtrFromNullString(ns) != s {
		t.Fatalf("unexpected null string behavior: %+v", ns)
	}
	empty := ""
	if NewNullString(&empty).Valid {
		t.Fatal("expected empty string to produce invalid sql.NullString")
	}
	if PtrFromNullString(sql.NullString{}) != nil {
		t.Fatal("expected nil pointer from invalid sql.NullString")
	}

	now := time.Now().UTC().Truncate(time.Second)
	nt := NewNullTime(&now)
	if !nt.Valid || !nt.Time.Equal(now) {
		t.Fatalf("unexpected null time behavior: %+v", nt)
	}
	tPtr := PtrFromNullTime(nt)
	if tPtr == nil || !tPtr.Equal(now) {
		t.Fatalf("unexpected time pointer from null: %v", tPtr)
	}
	if NewNullTime(nil).Valid || PtrFromNullTime(sql.NullTime{}) != nil {
		t.Fatal("expected invalid/null behavior for null time")
	}
}

