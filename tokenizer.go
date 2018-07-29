package main

import (
	"fmt"
	"strings"
	"unicode"
)

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

//TokenDebugPrint prints the contents of the token array in a standardized format
func TokenDebugPrint(tokens []Token) {
	for _, token := range tokens {
		fmt.Printf("  L%03d,C%03d-%03d %11s: [%s]\n", token.line, token.colStart, token.colEnd, token.tokType.toString(), token.value)
	}
}

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

func (tokType TokenType) toString() string {
	curTokType := "None"
	switch tokType {
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
