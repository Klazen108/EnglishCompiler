package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
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

	statements := Compile(dat)
	StatementDebugPrint(statements)

	ast := genAST(statements)

	state := ProgramState{
		identifiers: map[string]string{},
		types:       map[string]DataType{},
	}
	RunProgram(ast, state)
	fmt.Println("Program Complete! Result:")
	state.PrintState()
}

//RunProgram executes the program (represented as an Abstract Syntax Tree) against a backing program state.
func RunProgram(ast AST, state ProgramState) {
	for _, verb := range ast {
		verb.evaluate(state)
	}
}

func (s ProgramState) PrintState() {
	for k, v := range s.identifiers {
		fmt.Printf("(%6s)%s = %s\n", s.types[k].toString(), k, v)
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

func genAST(program []Statement) AST {
	var ast AST
	for stmtIndex, stmt := range program {
		if stmt[0].tokType != Word {
			panic(ParserError{stmtIndex: stmtIndex, tok: stmt[0], message: fmt.Sprintf("First token in statement must be a word, was %s", stmt[0].tokTypeString())})
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

//SetVerb modifies the program state by setting the value of the
//variable identified by alpha to the value returned by beta.
type SetVerb struct {
	alpha Identifier
	beta  Expression
}

//AddVerb modifies the program state by adding an expression's result
//to the variable identified by beta.
type AddVerb struct {
	alpha Expression
	beta  Identifier
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

func (v SetVerb) evaluate(state ProgramState) {
	state.setIdentifier(v.alpha, v.beta.evaluate(state), v.beta.dataType())
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
	state.setIdentifier(v.beta, s, dtNumber)
}

//Return: (type, exists?)
func (s ProgramState) getType(i Identifier) (DataType, bool) {
	t, ok := s.types[i.name]
	return t, ok
}

func (s ProgramState) getValue(i Identifier) (string, bool) {
	t, ok := s.identifiers[i.name]
	return t, ok
}

//An AST in the context of this interpreter is a list of verbs
//to be executed in sequence. More generally, an AST is an abstract
//syntax tree, representing a program source code in a more
//interpreter-friendly format.
type AST []Verb

//A Verb is an executable expression which returns no result, but
//may modify program (or global) state.
type Verb interface {
	evaluate(state ProgramState)
}

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

//The ProgramState encapsulates the complete state of the program.
//Identifiers identify memory in this structure
type ProgramState struct {
	identifiers map[string]string
	types       map[string]DataType
}

func (state ProgramState) setIdentifier(id Identifier, value string, dType DataType) {
	if vType, ok := state.types[id.name]; ok {
		if vType != dType {
			panic("Invalid datatype!")
		}
	} else {
		state.types[id.name] = dType
	}
	state.identifiers[id.name] = value
}

type DataType uint

const (
	dtAny    DataType = 0
	dtString DataType = 1
	dtNumber DataType = 2
)

func Compile(input []byte) []Statement {
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

//A Statement is a collection of tokens representing a single operation/action.
type Statement []Token

//A Token is an individual unit of code, ready to be lexically analyzed.
type Token struct {
	tokType  TokenType
	value    string
	line     uint
	colStart uint
	colEnd   uint
	purpose  PurposeType
}

/*TokenType represents a type of token*/
type TokenType uint

const (
	//None represents an unrecognized token
	None TokenType = 0
	//Punctuation represents separator tokens
	Punctuation TokenType = 1
	//Whitespace represents blank space tokens
	Whitespace TokenType = 2
	//Number represents numeric constant tokens
	Number TokenType = 3
	//Word represents variables, keywords, and character-based constant tokens
	Word TokenType = 4
)

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

//TokenDebugPrint prints the contents of the token array in a standardized format
func TokenDebugPrint(tokens []Token) {
	for _, token := range tokens {
		fmt.Printf("  L%03d,C%03d-%03d %11s: [%s]\n", token.line, token.colStart, token.colEnd, token.tokTypeString(), token.value)
	}
}

func StatementDebugPrint(statements []Statement) {
	for index, statement := range statements {
		fmt.Printf("Statement %03d\n", index)
		TokenDebugPrint(statement)
	}
}

//tokTypeString returns a human-readable interpretation of each token type
func (token Token) tokTypeString() string {
	curTokType := "None"
	switch token.tokType {
	case None:
		curTokType = "None"
		break
	case Punctuation:
		curTokType = "Punctuation"
		break
	case Whitespace:
		curTokType = "Whitespace"
		break
	case Number:
		curTokType = "Number"
		break
	case Word:
		curTokType = "Word"
		break
	}
	return curTokType
}

//split splits a token into two
//offset: 0-based index of the first char to split off for the right token
func (token Token) split(offset uint, leftType TokenType, rightType TokenType) [2]Token {
	return [2]Token{
		Token{
			tokType:  leftType,
			value:    token.value[0:offset],
			line:     token.line,
			colStart: token.colStart,
			colEnd:   token.colStart + offset - 1,
		},
		Token{
			tokType:  rightType,
			value:    token.value[offset:],
			line:     token.line,
			colStart: token.colStart + offset,
			colEnd:   token.colEnd},
	}
}

//Tokenize takes a unicode input and scans over it,
//determining token boundaries and their types.
func Tokenize(input []byte) []Token {
	curToken := ""
	var tokens []Token
	curTokenType := None
	var curLine uint = 1
	var curCol uint = 1
	var curColStart uint = 1
	for _, runeValue := range input {
		nextTokenType := curTokenType
		switch curTokenType {
		case None:
			if unicode.IsPunct(rune(runeValue)) {
				nextTokenType = Punctuation
			} else if unicode.IsSpace(rune(runeValue)) {
				nextTokenType = Whitespace
			} else if unicode.IsDigit(rune(runeValue)) {
				nextTokenType = Number
			} else if unicode.IsLetter(rune(runeValue)) {
				nextTokenType = Word
			}
			break
		case Punctuation:
			if unicode.IsSpace(rune(runeValue)) {
				nextTokenType = Whitespace
			} else if unicode.IsDigit(rune(runeValue)) {
				nextTokenType = Number
			} else if unicode.IsLetter(rune(runeValue)) {
				nextTokenType = Word
			}
			break
		case Whitespace:
			if unicode.IsPunct(rune(runeValue)) {
				nextTokenType = Punctuation
			} else if unicode.IsDigit(rune(runeValue)) || stringInSlice(rune(runeValue), []rune{rune('-'), rune('$')}) {
				nextTokenType = Number
			} else if unicode.IsLetter(rune(runeValue)) {
				nextTokenType = Word
			}
			break
		case Number:
			if unicode.IsPunct(rune(runeValue)) && !stringInSlice(rune(runeValue), []rune{rune(','), rune('.')}) {
				nextTokenType = Punctuation
			} else if unicode.IsSpace(rune(runeValue)) {
				nextTokenType = Whitespace
			}
			break
		case Word:
			if unicode.IsPunct(rune(runeValue)) {
				nextTokenType = Punctuation
			} else if unicode.IsSpace(rune(runeValue)) {
				nextTokenType = Whitespace
			}
			break
		}
		if nextTokenType == curTokenType {
			curToken += string(runeValue)
		} else {
			if curTokenType != None {
				tokens = append(tokens, Token{
					tokType:  curTokenType,
					value:    curToken,
					line:     curLine,
					colStart: curColStart,
					colEnd:   curCol - 1, //the last token ended on the previous character
				})
				curColStart = curCol
			}
			curToken = string(runeValue)
			curTokenType = nextTokenType
		}
		curCol++
		if rune(runeValue) == rune('\n') {
			curCol = 1
			curLine++
		}
	}
	tokens = append(tokens, Token{
		tokType:  curTokenType,
		value:    curToken,
		line:     curLine,
		colStart: curColStart,
		colEnd:   curCol - 1,
	})
	return tokens
}

//TokenizeStep2 performs the second tokenization step. Since there are some shared character sets
//between token types, it is necessary to check those types and split the tokens apart if necessary.
//For example, "0.0" is a valid number, but "0." should be two separate number and punctuation tokens.
func TokenizeStep2(tokens []Token) []Token {
	var newTokens []Token
	for _, token := range tokens {
		if token.tokType == Number && strings.HasSuffix(token.value, ".") {
			splitTokens := token.split(uint(len(token.value))-1, Number, Punctuation)
			newTokens = append(newTokens, splitTokens[:]...)
			continue
		}
		if token.tokType == Number && strings.HasSuffix(token.value, ",") {
			splitTokens := token.split(uint(len(token.value))-1, Number, Punctuation)
			newTokens = append(newTokens, splitTokens[:]...)
			continue
		}
		if token.tokType == Punctuation && strings.HasSuffix(token.value, ".") && strings.Contains(token.value, "\"") {
			splitTokens := token.split(uint(len(token.value))-1, Number, Punctuation)
			newTokens = append(newTokens, splitTokens[:]...)
			continue
		}
		newTokens = append(newTokens, token)
	}
	return newTokens
}
