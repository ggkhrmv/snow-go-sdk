package table

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyQueryField      = errors.New("query field cannot be empty")
	ErrInvalidQueryField    = errors.New("query field contains invalid characters")
	ErrEmptyQueryValue      = errors.New("query value cannot be empty")
	ErrEmptyQueryValues     = errors.New("query requires at least one value")
	ErrEmptyQueryOperator   = errors.New("query operator cannot be empty")
	ErrInvalidQueryLogical  = errors.New("logical operator cannot be the first query token")
	ErrDanglingQueryLogical = errors.New("query ends with a dangling logical operator")
)

// QueryBuilder composes a sysparm_query encoded query string.
//
// Example:
//
//	query, err := table.NewQueryBuilder().
//		Eq("active", true).
//		Or().
//		Eq("priority", 1).
//		Build()
type QueryBuilder struct {
	parts        []string
	nextOperator string
	err          error
}

// NewQueryBuilder creates a new encoded-query builder.
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		nextOperator: "^",
	}
}

func (b *QueryBuilder) Eq(field string, value any) *QueryBuilder {
	return b.addBinary(field, "=", value)
}
func (b *QueryBuilder) NotEq(field string, value any) *QueryBuilder {
	return b.addBinary(field, "!=", value)
}
func (b *QueryBuilder) GT(field string, value any) *QueryBuilder {
	return b.addBinary(field, ">", value)
}
func (b *QueryBuilder) GTE(field string, value any) *QueryBuilder {
	return b.addBinary(field, ">=", value)
}
func (b *QueryBuilder) LT(field string, value any) *QueryBuilder {
	return b.addBinary(field, "<", value)
}
func (b *QueryBuilder) LTE(field string, value any) *QueryBuilder {
	return b.addBinary(field, "<=", value)
}
func (b *QueryBuilder) Contains(field string, value any) *QueryBuilder {
	return b.addBinary(field, "LIKE", value)
}
func (b *QueryBuilder) NotContains(field string, value any) *QueryBuilder {
	return b.addBinary(field, "NOT LIKE", value)
}
func (b *QueryBuilder) StartsWith(field string, value any) *QueryBuilder {
	return b.addBinary(field, "STARTSWITH", value)
}
func (b *QueryBuilder) EndsWith(field string, value any) *QueryBuilder {
	return b.addBinary(field, "ENDSWITH", value)
}
func (b *QueryBuilder) IsEmpty(field string) *QueryBuilder    { return b.addUnary(field, "ISEMPTY") }
func (b *QueryBuilder) IsNotEmpty(field string) *QueryBuilder { return b.addUnary(field, "ISNOTEMPTY") }

// Op appends a condition with a custom ServiceNow operator.
func (b *QueryBuilder) Op(field, operator string, value any) *QueryBuilder {
	operator = strings.TrimSpace(operator)
	if operator == "" {
		b.setErr(ErrEmptyQueryOperator)
		return b
	}

	return b.addBinary(field, operator, value)
}

func (b *QueryBuilder) In(field string, values ...any) *QueryBuilder {
	return b.addList(field, "IN", values...)
}

func (b *QueryBuilder) NotIn(field string, values ...any) *QueryBuilder {
	return b.addList(field, "NOT IN", values...)
}

// And sets AND conjunction for the next condition.
func (b *QueryBuilder) And() *QueryBuilder {
	return b.setLogical("^")
}

// Or sets OR conjunction for the next condition.
func (b *QueryBuilder) Or() *QueryBuilder {
	return b.setLogical("^OR")
}

// NewQuery starts a new encoded-query group (^NQ).
func (b *QueryBuilder) NewQuery() *QueryBuilder {
	return b.setLogical("^NQ")
}

// Build returns the encoded query string.
func (b *QueryBuilder) Build() (string, error) {
	if b == nil {
		return "", nil
	}
	if b.err != nil {
		return "", b.err
	}
	if len(b.parts) == 0 {
		return "", nil
	}
	if b.nextOperator != "^" {
		return "", ErrDanglingQueryLogical
	}

	return strings.Join(b.parts, ""), nil
}

// String returns the built encoded query and discards build errors.
func (b *QueryBuilder) String() string {
	query, _ := b.Build()
	return query
}

// Err returns the first builder error.
func (b *QueryBuilder) Err() error {
	if b == nil {
		return nil
	}
	if b.err != nil {
		return b.err
	}
	if len(b.parts) > 0 && b.nextOperator != "^" {
		return ErrDanglingQueryLogical
	}

	return nil
}

func (b *QueryBuilder) setLogical(op string) *QueryBuilder {
	if b == nil || b.err != nil {
		return b
	}
	if len(b.parts) == 0 {
		b.setErr(ErrInvalidQueryLogical)
		return b
	}

	b.nextOperator = op
	return b
}

func (b *QueryBuilder) addBinary(field, operator string, value any) *QueryBuilder {
	field, ok := b.validField(field)
	if !ok {
		return b
	}

	valueString, err := queryValue(value)
	if err != nil {
		b.setErr(err)
		return b
	}

	return b.addCondition(field + operator + valueString)
}

func (b *QueryBuilder) addUnary(field, operator string) *QueryBuilder {
	field, ok := b.validField(field)
	if !ok {
		return b
	}

	return b.addCondition(field + operator)
}

func (b *QueryBuilder) addList(field, operator string, values ...any) *QueryBuilder {
	field, ok := b.validField(field)
	if !ok {
		return b
	}
	if len(values) == 0 {
		b.setErr(ErrEmptyQueryValues)
		return b
	}

	formatted := make([]string, 0, len(values))
	for _, value := range values {
		valueString, err := queryValue(value)
		if err != nil {
			b.setErr(err)
			return b
		}
		formatted = append(formatted, valueString)
	}

	return b.addCondition(field + operator + strings.Join(formatted, ","))
}

func (b *QueryBuilder) validField(field string) (string, bool) {
	field = strings.TrimSpace(field)
	if field == "" {
		b.setErr(ErrEmptyQueryField)
		return "", false
	}
	if strings.Contains(field, "^") {
		b.setErr(ErrInvalidQueryField)
		return "", false
	}

	return field, true
}

func (b *QueryBuilder) addCondition(condition string) *QueryBuilder {
	if b == nil || b.err != nil {
		return b
	}
	if condition == "" {
		b.setErr(ErrEmptyQueryValue)
		return b
	}

	if len(b.parts) == 0 {
		b.parts = append(b.parts, condition)
		b.nextOperator = "^"
		return b
	}

	b.parts = append(b.parts, b.nextOperator+condition)
	b.nextOperator = "^"

	return b
}

func (b *QueryBuilder) setErr(err error) {
	if b == nil || b.err != nil || err == nil {
		return
	}
	b.err = err
}

func queryValue(value any) (string, error) {
	if value == nil {
		return "", ErrEmptyQueryValue
	}

	raw := fmt.Sprint(value)
	if raw == "" {
		return "", ErrEmptyQueryValue
	}

	return strings.ReplaceAll(raw, "^", "^^"), nil
}
