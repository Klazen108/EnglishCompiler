package main

import "strconv"

//An Expression is an evaluatable expression which returns a result
//and has a data type.
type Expression interface {
	dataType() DataType
	evaluate(state ProgramState) string
	toString() string
}

type IdentifierExpression struct {
	id Identifier
}

func (c IdentifierExpression) evaluate(state ProgramState) string {
	val, ok := state.getValue(c.id)
	if !ok {
		panic("invalid identifier: " + c.id.name)
	}
	return val
}

func (c IdentifierExpression) dataType() DataType {
	return dtNumber //TODO: idk?
}

func (e IdentifierExpression) toString() string {
	return e.id.name
}

//ConstNumExpression is an expression which evaluates to the
//same constant numeric expression, always
type ConstNumExpression struct {
	value string
}

//ConstNumExpression.dataType always returns dtNumber, since it
//is constant and numeric
func (c ConstNumExpression) dataType() DataType {
	return dtNumber
}

func (c ConstNumExpression) evaluate(state ProgramState) string {
	return c.value
}

func (c ConstNumExpression) toString() string {
	return c.value
}

type GreaterThanExpression struct {
	lhs Expression
	rhs Expression
}

func (e GreaterThanExpression) evaluate(state ProgramState) string {
	if e.lhs.dataType() != dtNumber {
		panic("Invalid datatype for LHS: " + e.lhs.dataType().toString())
	}
	if e.rhs.dataType() != dtNumber {
		panic("Invalid datatype for RHS: " + e.rhs.dataType().toString())
	}

	sLhs := e.lhs.evaluate(state)
	iLhs, err := strconv.Atoi(sLhs)
	if err != nil {
		panic("Unable to parse LHS, despite having a numeric datatype: " + sLhs)
	}

	sRhs := e.rhs.evaluate(state)
	iRhs, err := strconv.Atoi(sRhs)
	if err != nil {
		panic("Unable to parse RHS, despite having a numeric datatype: " + sRhs)
	}

	if iLhs > iRhs {
		return "1"
	}
	return "0"
}

func (e GreaterThanExpression) dataType() DataType {
	return dtNumber
}

func (e GreaterThanExpression) toString() string {
	return e.lhs.toString() + " > " + e.rhs.toString()
}
