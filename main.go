package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

	tokens := Tokenize(dat)
	TokenDebugPrint(tokens)
	//fmt.Printf("%v", tokens)
}

type Token struct {
	tokType  TokenType
	value    string
	line     uint
	colStart uint
	colEnd   uint
}

/*TokenType represents a type of token*/
type TokenType uint

const (
	None        TokenType = 0
	Punctuation TokenType = 1
	Whitespace  TokenType = 2
	Number      TokenType = 3
	Word        TokenType = 4
)

func stringInSlice(a rune, list []rune) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func TokenDebugPrint(tokens []Token) {
	for _, token := range tokens {
		fmt.Printf("L%03d,C%03d-%03d %11s: [%s]\n", token.line, token.colStart, token.colEnd, token.tokTypeString(), token.value)
	}
}

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
			} else if unicode.IsLetter(rune(runeValue)) {
				nextTokenType = Word
			}
			break
		case Word:
			if unicode.IsPunct(rune(runeValue)) {
				nextTokenType = Punctuation
			} else if unicode.IsSpace(rune(runeValue)) {
				nextTokenType = Whitespace
			} else if unicode.IsDigit(rune(runeValue)) {
				nextTokenType = Number
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
		colEnd:   curCol,
	})
	return tokens
}
