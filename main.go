package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("English Compiler - Klazen108, 2018")

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: eng src.eng")
		return
	}
	file := args[1]
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	fmt.Println("Program Input:")
	fmt.Println(string(dat))

	fmt.Println("Compiling...")
	statements := Compile(dat)
	//StatementDebugPrint(statements)

	ast := GenAST(statements)

	fmt.Println("Executing...")
	state := ProgramState{
		identifiers: map[string]string{},
		types:       map[string]DataType{},
	}
	RunProgram(ast, state)
	fmt.Println("Final Program State:")
	state.PrintState()
}

//RunProgram executes the program (represented as an Abstract Syntax Tree) against a backing program state.
func RunProgram(ast AST, state ProgramState) {
	for _, verb := range ast {
		verb.evaluate(state)
	}
}

type PurposeType uint

const (
	ptKeyword    PurposeType = 1
	ptVerb       PurposeType = 2
	ptIdentifier PurposeType = 3
	ptConstant   PurposeType = 4
)

func stringInSlice(a rune, list []rune) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
