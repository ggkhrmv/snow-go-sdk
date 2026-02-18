package table

import (
	"errors"
	"testing"
)

func TestQueryBuilderBuild(t *testing.T) {
	query, err := NewQueryBuilder().
		Eq("active", true).
		GT("priority", 2).
		Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	const want = "active=true^priority>2"
	if query != want {
		t.Fatalf("Build() = %q, want %q", query, want)
	}
}

func TestQueryBuilderLogicalOperators(t *testing.T) {
	query, err := NewQueryBuilder().
		Eq("state", 1).
		Or().
		Eq("state", 2).
		NewQuery().
		IsNotEmpty("assigned_to").
		Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	const want = "state=1^ORstate=2^NQassigned_toISNOTEMPTY"
	if query != want {
		t.Fatalf("Build() = %q, want %q", query, want)
	}
}

func TestQueryBuilderIn(t *testing.T) {
	query, err := NewQueryBuilder().
		In("priority", 1, 2, 3).
		Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	const want = "priorityIN1,2,3"
	if query != want {
		t.Fatalf("Build() = %q, want %q", query, want)
	}
}

func TestQueryBuilderEscapesCaret(t *testing.T) {
	query, err := NewQueryBuilder().
		Eq("short_description", "foo^bar").
		Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	const want = "short_description=foo^^bar"
	if query != want {
		t.Fatalf("Build() = %q, want %q", query, want)
	}
}

func TestQueryBuilderEmptyField(t *testing.T) {
	_, err := NewQueryBuilder().
		Eq("", "x").
		Build()
	if !errors.Is(err, ErrEmptyQueryField) {
		t.Fatalf("Build() error = %v, want %v", err, ErrEmptyQueryField)
	}
}

func TestQueryBuilderDanglingLogicalOperator(t *testing.T) {
	_, err := NewQueryBuilder().
		Eq("active", true).
		Or().
		Build()
	if !errors.Is(err, ErrDanglingQueryLogical) {
		t.Fatalf("Build() error = %v, want %v", err, ErrDanglingQueryLogical)
	}
}

func TestQueryBuilderLogicalCannotBeFirst(t *testing.T) {
	_, err := NewQueryBuilder().
		Or().
		Eq("active", true).
		Build()
	if !errors.Is(err, ErrInvalidQueryLogical) {
		t.Fatalf("Build() error = %v, want %v", err, ErrInvalidQueryLogical)
	}
}

func TestQueryBuilderInRequiresValues(t *testing.T) {
	_, err := NewQueryBuilder().
		In("priority").
		Build()
	if !errors.Is(err, ErrEmptyQueryValues) {
		t.Fatalf("Build() error = %v, want %v", err, ErrEmptyQueryValues)
	}
}
