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
	for _, stmt := range program {
		v, err := ParseVerb(stmt)
		if err != nil { //TODO
			//panic(ParserError{stmtIndex: stmtIndex, tok: stmt[0], message: fmt.Sprintf("First token in statement must be a word, was %s", stmt[0].tokType.toString())})
			panic(err)
		}
		ast = append(ast, v)
	}
	return ast
}

func ParseVerb(stmt []Token) (Verb, error) {
	operation := strings.ToLower(stmt[0].value)

	if stmt[0].tokType != Word {
		return nil, nil //TODO
	}

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
		return v, nil
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
		return v, nil
	case "display":
		v := DisplayVerb{
			alpha: IdentifierExpression{
				id: Identifier{
					name:  stmt[1].value,
					dType: dtString,
				},
			},
		}
		return v, nil
	case "if":
		splitTokens, ok := SplitTokens("then", stmt[1:])
		if !ok {
			panic("Couldn't find a then in an if statemet")
		}
		sube, err := ParseExpression(splitTokens[0])
		if err != nil {
			return nil, err
		}
		subv, err := ParseVerb(splitTokens[1])
		if err != nil {
			return nil, err
		}
		if subv == nil {
			return nil, ParserError{message: "Couldn't find action for if statement", stmtIndex: 0, tok: stmt[0]}
		}
		v := IfVerb{
			predicate: sube,
			action:    subv,
		}
		return v, nil
	case "otherwise":
		subv, err := ParseVerb(stmt[1:])
		if err != nil {
			return nil, err
		}
		v := OtherwiseVerb{
			action: subv,
		}
		return v, nil
	default:
		return nil, nil //TODO
	}
}

type Matcher struct {
	entries []MatchEntry
}

func (m Matcher) reset() {
	for i, _ := range m.entries {
		m.entries[i].matched = false
	}
}

func (m Matcher) matched() bool {
	for _, entry := range m.entries {
		if !entry.matched {
			return false
		}
	}
	return true
}

func (m Matcher) checkMatch(word string) {
	for i, entry := range m.entries {
		if entry.matched {
			continue
		}
		if entry.word == word {
			m.entries[i].matched = true
		} else {
			m.reset()
		}
		break
	}
}

type MatchEntry struct {
	word    string
	matched bool
}

type MatcherError struct {
}

func (m MatcherError) Error() string {
	return "matcher error"
}

func ParseExpression(tokens []Token) (Expression, error) {
	//fmt.Printf("%v\n", tokens)
	matchList := []Matcher{
		Matcher{
			entries: []MatchEntry{
				MatchEntry{word: "is", matched: false},
				MatchEntry{word: "greater", matched: false},
				MatchEntry{word: "than", matched: false},
			},
		},
	}
	for _, token := range tokens {
		for _, m := range matchList {
			m.checkMatch(token.value)
			if m.matched() {
				//found one, parse it
				//TODO: get section before and after matched entry, and parse expressions out of those
				return GreaterThanExpression{
					lhs: IdentifierExpression{id: Identifier{dType: dtNumber, name: "value"}},
					rhs: ConstNumExpression{value: "50"},
				}, nil
			}
		}
	}
	return nil, MatcherError{}
}

func SplitTokens(needle string, tokens []Token) ([2][]Token, bool) {
	//fmt.Printf("%v - %s\n", tokens, needle)
	for i, token := range tokens {
		if strings.ToLower(token.value) == needle {
			//fmt.Printf("found at %d\n", i)
			return [2][]Token{tokens[0:i], tokens[i+1:]}, true
		}
	}
	return [2][]Token{}, false
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
