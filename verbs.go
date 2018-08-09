package main

import (
	"fmt"
	"strconv"
)

//A Verb is an executable expression which returns no result, but
//may modify program (or global) state.
type Verb interface {
	evaluate(state ProgramState)
	toString() string
}

//SetVerb modifies the program state by setting the value of the
//variable identified by alpha to the value returned by beta.
type SetVerb struct {
	alpha Identifier
	beta  Expression
}

func (v SetVerb) evaluate(state ProgramState) {
	if err := state.SetValue(v.alpha, v.beta.evaluate(state), v.beta.dataType()); err != nil {
		panic(err)
	}
}

func (v SetVerb) toString() string {
	return v.alpha.name + " = " + v.beta.toString()
}

//AddVerb modifies the program state by adding an expression's result
//to the variable identified by beta.
type AddVerb struct {
	alpha Expression
	beta  Identifier
}

func (v AddVerb) evaluate(state ProgramState) {
	dType, exists := state.getType(v.beta)
	if !exists {
		panic("Uninitialized identifier " + v.beta.name)
	}

	if dType != dtNumber {
		panic("Invalid datatype for identifier " + v.beta.name + "! Expected number, got " + dType.toString())
	}

	if v.alpha.dataType() != dtNumber {
		panic("Invalid datatype for expression! Expected number, got " + v.alpha.dataType().toString())
	}

	sValue, _ := state.getValue(v.beta)
	i, err := strconv.Atoi(sValue)
	if err != nil {
		panic("Unable to parse integer, despite having a numeric datatype: " + sValue)
	}
	sAddend := v.alpha.evaluate(state)
	iAddend, err := strconv.Atoi(sAddend)
	if err != nil {
		panic("Unable to parse integer, despite having a numeric datatype: " + sAddend)
	}
	i += iAddend
	s := strconv.Itoa(i)
	if err = state.SetValue(v.beta, s, dtNumber); err != nil {
		panic(err)
	}
}

func (v AddVerb) toString() string {
	return v.beta.name + "+=" + v.alpha.toString()
}

//DisplayVerb prints the value of its expression argument to the console
type DisplayVerb struct {
	alpha Expression
}

func (v DisplayVerb) evaluate(state ProgramState) {
	fmt.Println(v.alpha.evaluate(state))
}

func (v DisplayVerb) toString() string {
	return "display " + v.alpha.toString()
}

type IfVerb struct {
	predicate Expression
	action    Verb
}

func (v IfVerb) evaluate(state ProgramState) {
	if v.predicate.evaluate(state) == "1" {
		fmt.Printf("%s was true: doing %s", v.predicate.toString(), v.action.toString())
		v.action.evaluate(state)
	} else {
		state.FlagElse(true)
	}
}

func (v IfVerb) toString() string {
	return "if " + v.predicate.toString() + " then " + v.action.toString()
}

type OtherwiseVerb struct {
	action Verb
}

func (v OtherwiseVerb) evaluate(state ProgramState) {
	if state.IsFlaggedElse() {
		v.action.evaluate(state)
		state.FlagElse(false)
	}
}

func (v OtherwiseVerb) toString() string {
	return "otherwise " + v.action.toString()
}
