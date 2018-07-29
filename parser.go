package main

import (
	"fmt"
	"strings"
)

//A Statement is a collection of tokens representing a single operation/action.
type Statement []Token

//A program is a collection of statements. To execute the program, it first
//needs to be converted into an AST - see GenAST
type Program []Statement

func StatementDebugPrint(statements Program) {
	for index, statement := range statements {
		fmt.Printf("Statement %03d\n", index)
		TokenDebugPrint(statement)
	}
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

func (d DataType) toString() string {
	switch d {
	case dtAny:
		return "Any"
	case dtNumber:
		return "Number"
	case dtString:
		return "String"
	default:
		return "Unknown"
	}
}

//An AST in the context of this interpreter is a list of verbs
//to be executed in sequence. More generally, an AST is an abstract
//syntax tree, representing a program source code in a more
//interpreter-friendly format.
type AST []Verb

type DataType uint

const (
	dtAny    DataType = 0
	dtString DataType = 1
	dtNumber DataType = 2
)

//An Expression is an evaluatable expression which returns a result
//and has a data type.
type Expression interface {
	dataType() DataType
	evaluate(state ProgramState) string
}

//An Identifier is a "pointer" to some memory location in program state,
//referenced by name and with an associated data type.
type Identifier struct {
	name  string
	dType DataType
}

//Compile takes an input document and converts it to a Program. Before
//execution, it needs to be converted to an AST - see GenAST
func Compile(input []byte) Program {
	tokens := Tokenize(input)
	tokens = TokenizeStep2(tokens)
	//remove whitespace tokens - they're unnecessary past this point
	var newTokens []Token
	for _, token := range tokens {
		if token.tokType != Whitespace {
			newTokens = append(newTokens, token)
		}
	}
	tokens = newTokens

	//TokenDebugPrint(tokens)

	//Split the stream of tokens into statements
	var curStatement Statement
	var statements []Statement
	for _, token := range tokens {
		if token.tokType == Punctuation && token.value == "." {
			statements = append(statements, curStatement)
			curStatement = Statement{}
		} else {
			//no need to carry the period past here
			curStatement = append(curStatement, token)
		}
	}
	return statements
}

//GenAST takes a collection of program statements, and creates an AST
//which is executable.
func GenAST(program []Statement) AST {
	var ast AST
	for stmtIndex, stmt := range program {
		if stmt[0].tokType != Word {
			panic(ParserError{stmtIndex: stmtIndex, tok: stmt[0], message: fmt.Sprintf("First token in statement must be a word, was %s", stmt[0].tokType.toString())})
		}

		operation := strings.ToLower(stmt[0].value)
		switch operation {
		case "set":
			//TODO: parse and generate the alpha and beta dynamically, don't assume it will always be "set alpha to beta"
			/*
				i := 1
				for i = 1; i < len(stmt); i++ {
					if strings.ToLower(stmt[i].value) == "to" {
						break
					}
				}
				funcSet(stmt[1:i-1], stmt[i+1:], state)
			*/
			v := SetVerb{
				alpha: Identifier{
					name:  stmt[1].value,
					dType: dtAny,
				},
				beta: ConstNumExpression{
					value: stmt[3].value,
				},
			}
			ast = append(ast, v)
			break
		case "add":
			v := AddVerb{
				alpha: ConstNumExpression{
					value: stmt[1].value,
				},
				beta: Identifier{
					name:  stmt[3].value,
					dType: dtAny,
				},
			}
			ast = append(ast, v)
			break
		}
	}
	return ast
}

//ParserError represents a failure in parsing.
type ParserError struct {
	stmtIndex int
	tok       Token
	message   string
}

//ParserError.Error displays the statement, line and column numbers (for debugging) as well as an error message explaining
//what failed to parse.
func (p ParserError) Error() string {
	return fmt.Sprintf("Parser Failure: Statement %d (Line %d Column %d) %s", p.stmtIndex, p.tok.line, p.tok.colStart, p.message)
}
