package main

import "testing"

func TestGreaterThan(t *testing.T) {
	s := ProgramState{
		identifiers: map[string]string{},
		types:       map[string]DataType{},
	}

	tables := []struct {
		lhs      Expression
		rhs      Expression
		expected string
	}{
		{
			lhs:      ConstNumExpression{value: "5"},
			rhs:      ConstNumExpression{value: "10"},
			expected: "0",
		},
		{
			lhs:      ConstNumExpression{value: "10"},
			rhs:      ConstNumExpression{value: "5"},
			expected: "1",
		},
		{
			lhs:      ConstNumExpression{value: "5"},
			rhs:      ConstNumExpression{value: "5"},
			expected: "0",
		},
	}

	for _, table := range tables {
		e := GreaterThanExpression{
			lhs: table.lhs,
			rhs: table.rhs,
		}
		actual := e.evaluate(s)
		if actual != table.expected {
			t.Errorf("Expected %s got %s", table.expected, actual)
		}
	}
}
