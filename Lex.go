package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

const (
	// TkIllegal Special tokens
	TkIllegal int8 = 0 // __nil__
	TkEof     int8 = 1 // \0
	TkComment int8 = 2 // "//" || "/**/"

	// TkIdent Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	TkIdent  int8 = 3 // main
	TkInt    int8 = 4 // 12345
	TkFloat  int8 = 5 // 123.45
	TkImag   int8 = 6 // 123.45i
	TkChar   int8 = 7 // 'a'
	TkString int8 = 8 // "abc"

	// TkAdd Operators and delimiters
	TkAdd          int8 = 9  // +
	TkSub          int8 = 10 // -
	TkMul          int8 = 11 // *
	TkQuo          int8 = 12 // /
	TkRem          int8 = 13 // %
	TkAnd          int8 = 14 // &
	TkOr           int8 = 15 // |
	TkXor          int8 = 16 // ^
	TkShl          int8 = 17 // <<
	TkShr          int8 = 18 // >>
	TkAndNot       int8 = 19 // &^
	TkAddAssign    int8 = 20 // +=
	TkSubAssign    int8 = 21 // -=
	TkMulAssign    int8 = 22 // *=
	TkQuoAssign    int8 = 23 // /=
	TkRemAssign    int8 = 24 // %=
	TkAndAssign    int8 = 25 // &=
	TkOrAssign     int8 = 26 // |=
	TkXorAssign    int8 = 27 // ^=
	TkShlAssign    int8 = 28 // <<=
	TkShrAssign    int8 = 29 // >>=
	TkAndNotAssign int8 = 30 // &^=
	TkLand         int8 = 31 // &&
	TkLor          int8 = 32 // ||
	TkArrow        int8 = 33 // <-
	TkInc          int8 = 34 // ++
	TkDec          int8 = 35 // --
	TkEql          int8 = 36 // ==
	TkLss          int8 = 37 // <
	TkGtr          int8 = 38 // >
	TkAssign       int8 = 39 // =
	TkNot          int8 = 40 // !
	TkNeq          int8 = 41 // !=
	TkLeq          int8 = 42 // <=
	TkGeq          int8 = 43 // >=
	TkDefine       int8 = 44 // :=
	TkEllipsis     int8 = 45 // ...
	TkLParen       int8 = 46 // (
	TkLBrack       int8 = 47 // [
	TkLBrace       int8 = 48 // {
	TkComma        int8 = 49 // ,
	TkPeriod       int8 = 50 // .
	TkRParen       int8 = 51 // )
	TkRBrack       int8 = 52 // ]
	TkRBrace       int8 = 53 // }
	TkSemicolon    int8 = 54 // ;
	TkColon        int8 = 55 // :

	// TkBreakKw Keywords
	TkBreakKw       int8 = 56
	TkCaseKw        int8 = 57
	TkChanKw        int8 = 58
	TkConstKw       int8 = 59
	TkContinueKw    int8 = 60
	TkDefaultKw     int8 = 61
	TkDeferKw       int8 = 62
	TkElseKw        int8 = 63
	TkFallthroughKw int8 = 64
	TkForKw         int8 = 65
	TkFuncKw        int8 = 66
	TkGoKw          int8 = 67
	TkGotoKw        int8 = 68
	TkIfKw          int8 = 69
	TkImportKw      int8 = 70
	TkInterfaceKw   int8 = 71
	TkMapKw         int8 = 72
	TkPackageKw     int8 = 73
	TkRangeKw       int8 = 74
	TkReturnKw      int8 = 75
	TkSelectKw      int8 = 76
	TkStructKw      int8 = 77
	TkSwitchKw      int8 = 78
	TkTypeKw        int8 = 79
	TkVarKw         int8 = 80
	TkCleanKw       int8 = 81
	TkWhiteSpace    int8 = 82
	TkNewKw         int8 = 83
	TkIntKw         int8 = 84
	TkInt8Kw        int8 = 85
	TkInt16Kw       int8 = 86
	TkInt32Kw       int8 = 87
	TkInt64Kw       int8 = 88
	TkUIntKw        int8 = 89
	TkUInt8Kw       int8 = 90
	TkUInt16Kw      int8 = 91
	TkUInt32Kw      int8 = 92
	TKUint64Kw      int8 = 93
	TkFloatKw       int8 = 94
	TkFloat32Kw     int8 = 95
	TkFloat64Kw     int8 = 96
	TkStringKw      int8 = 97
	TkEnumKw        int8 = 98
	TkNilKw         int8 = 99
	TkAnyKw         int8 = 100

	// TkAbreviation Misc
	TkAbreviation int8 = 101 // /** This is an abreviation */
	TkVersionNum  int8 = 102
	TkSqrInc      int8 = 103

	TkByteKw int8 = 104
	TkBoolKw int8 = 105
)

var tokenTypeStrings []string = []string{
	"__tk_Illegal",
	"__tk_Eof",
	"__tk_Comment",
	"__tk_Ident",
	"__tk_Int",
	"__tk_Float",
	"__tk_Imag",
	"__tk_Char",
	"__tk_String",
	"__tk_Add",
	"__tk_Sub",
	"__tk_Mul",
	"__tk_Quo",
	"__tk_Rem",
	"__tk_And",
	"__tk_Or",
	"__tk_Xor",
	"__tk_Shl",
	"__tk_Shr",
	"__tk_AndNot",
	"__tk_AddAssign",
	"__tk_SubAssign",
	"__tk_MulAssign",
	"__tk_QuoAssign",
	"__tk_RemAssign",
	"__tk_AndAssign",
	"__tk_OrAssign",
	"__tk_XorAssign",
	"__tk_ShlAssign",
	"__tk_ShrAssign",
	"__tk_AndNotAssign",
	"__tk_Land",
	"__tk_Lor",
	"__tk_Arrow",
	"__tk_Inc",
	"__tk_Dec",
	"__tk_Eql",
	"__tk_Lss",
	"__tk_Gtr",
	"__tk_Assign",
	"__tk_Not",
	"__tk_Neq",
	"__tk_Leq",
	"__tk_Geq",
	"__tk_Define",
	"__tk_Ellipsis",
	"__tk_LParen",
	"__tk_LBrack",
	"__tk_LBrace",
	"__tk_Comma",
	"__tk_Period",
	"__tk_RParen",
	"__tk_RBrack",
	"__tk_RBrace",
	"__tk_Semicolon",
	"__tk_Colon",
	"break",
	"case",
	"chan",
	"const",
	"continue",
	"default",
	"defer",
	"else",
	"fallthrough",
	"for",
	"func",
	"go",
	"goto",
	"if",
	"import",
	"interface",
	"map",
	"package",
	"range",
	"return",
	"select",
	"struct",
	"switch",
	"type",
	"var",
	"clean",
	"__tk_WhiteSpace",
	"new",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"float",
	"float32",
	"float64",
	"string",
	"enum",
	"nil",
	"any",
	"__tk_Abreviation",
	"__tk_Version_Number",
	"__tk_Square_Increment",
	"byte",
	"bool",
}

type Token struct {
	tokenType int8
	lineno    uint
	sdata     string
	idata     int
	fdata     float64
}

func (gt Token) toString() string {
	return fmt.Sprintf("(%s, ln: %d, \"%s\", %d, %f)", tokenTypeStrings[gt.tokenType], gt.lineno, gt.sdata, gt.idata, gt.fdata)
}

func lex(txt string) ([]Token, []error) {
	data := []byte(txt)
	var tokens []Token
	var errs []error
	var ln uint = 1
	var i uint = 0

	data = append(data, '\n')

	for {
		if int(i) >= len(data)-1 {
			break
		}

		token := Token{}
		token.lineno = ln

		c := data[i]

		if unicode.IsLetter(rune(c)) || c == '_' {
			token.tokenType = TkIdent

			for unicode.IsLetter(rune(c)) || unicode.IsNumber(rune(c)) || c == '_' {
				token.sdata += string(c)

				i++
				c = data[i]

				if int(i) > len(data)-1 {
					errs = append(errs, errors.New("Identifier led to __eof__."+token.toString()))
				}
			}

			for _i := 56; _i < len(tokenTypeStrings); _i++ {
				if tokenTypeStrings[_i] == token.sdata {
					token.tokenType = int8(_i)
					break
				}
			}
		} else if unicode.IsNumber(rune(c)) {
			token.tokenType = TkInt
			var floatStr string //

			for unicode.IsDigit(rune(c)) || c == '\'' || c == '.' {
				if c == '\'' {
					i++
					c = data[i]
					continue
				} else if c == '.' {
					if token.tokenType == TkFloat || token.tokenType == TkVersionNum {
						token.tokenType = TkVersionNum
						floatStr += "."
					} else {
						token.tokenType = TkFloat
						floatStr += "."
					}
				} else {
					floatStr += string(c)
				}

				i++
				c = data[i]

				if int(i) > len(data)-1 {
					errs = append(errs, errors.New("Float/int literal led to __eof__."+token.toString()))
				}
			}

			if token.tokenType == TkFloat {
				fdata, err := strconv.ParseFloat(floatStr, 64)
				if err != nil {
					errs = append(errs, err)
				} else {
					token.fdata = fdata
				}
			} else if token.tokenType == TkVersionNum {
				token.sdata = floatStr
			} else {
				idata, err := strconv.ParseInt(floatStr, 10, 64)
				if err != nil {
					errs = append(errs, err)
				} else {
					token.idata = int(idata)
				}
			}
		} else {
			switch string(c) {
			case "\r":
				i++
				c = data[i]
				fallthrough
			case "\n":
				token.tokenType = TkWhiteSpace
				token.sdata += "__nil__"
				for string(c) == "\n" || string(c) == "\r" {
					if string(c) == "\n" {
						ln++
					}
					i++

					if int(i) > len(data)-1 {
						break
					}

					c = data[i]
				}

			case " ":
				fallthrough
			case "\t":
				i++

			case ".":
				if unicode.IsNumber(rune(data[i+1])) {
					token.tokenType = TkFloat
					floatStr := "0."

					i++
					c = data[i]

					for unicode.IsDigit(rune(c)) || c == '\'' || c == '.' {
						if c == '\'' {
							i++
							c = data[i]
							continue
						} else if c == '.' {
							token.tokenType = TkVersionNum
							floatStr += "."
						} else {
							floatStr += string(c)
						}

						i++
						c = data[i]

						if int(i) > len(data)-1 {
							errs = append(errs, errors.New("Float/int literal led to __eof__."+token.toString()))
						}
					}

					if token.tokenType == TkFloat {
						fdata, err := strconv.ParseFloat(floatStr, 64)
						if err != nil {
							errs = append(errs, err)
						} else {
							token.fdata = fdata
						}
					} else {
						token.sdata = floatStr
					}
				} else if data[i+1] == '.' && data[i+2] == '.' {
					i += 3
					token.tokenType = TkEllipsis
					token.sdata = "..."
				} else {
					token.tokenType = TkPeriod
					token.sdata = "."
					i++
				}

			case "/":
				if data[i+1] == '/' {
					i += 2
					c = data[i]
					token.tokenType = TkComment
					for c != '\n' {
						if string(c) != "\r" {
							token.sdata += string(c)
						}
						i++
						c = data[i]
					}
				} else if data[i+1] == '*' {
					i++
					if data[i+1] == '*' {
						i += 2
						c = data[i]
						token.tokenType = TkAbreviation

						for {
							if data[i+1] == '*' {
								if data[i+2] == '/' {
									i += 3
									break
								}
							}

							if string(c) != "\r" {
								token.sdata += string(c)
							}
							i++
							c = data[i]

							if int(i) > len(data)-1 {
								errs = append(errs, errors.New("Abriviation had no ending comment symbol."+token.toString()))
							}
						}
					} else {
						i++
						c = data[i]
						token.tokenType = TkComment

						for {
							if data[i+1] == '*' {
								if data[i+2] == '/' {
									i += 3
									break
								}
							}

							if string(c) != "\r" {
								token.sdata += string(c)
							}
							i++
							c = data[i]

							if int(i) > len(data)-1 {
								errs = append(errs, errors.New("Multi-line comment had no ending comment symbol."+token.toString()))
							}
						}
					}
				} else if data[i+1] == '=' {
					i += 2
					token.tokenType = TkQuoAssign
					token.sdata = "/="
				} else {
					i++
					token.tokenType = TkQuo
					token.sdata = "/"
				}

			case "\"":
				token.tokenType = TkString

				i++
				c = data[i]

				for c != '"' {
					token.sdata += string(c)
					i++
					c = data[i]

					if int(i) > len(data)-1 {
						errs = append(errs, errors.New("String literal led to __eof__"+token.toString()))
					}
				}

				i++

			case ":":
				if data[i+1] == '=' {
					i += 2
					token.tokenType = TkDefine
					token.sdata = ":="
				} else {
					i++
					token.tokenType = TkColon
					token.sdata = ":"
				}

			case ";":
				i++
				token.tokenType = TkSemicolon
				token.sdata = ";"

			case "*":
				if data[i+1] == '=' {
					i += 2
					token.tokenType = TkMulAssign
					token.sdata = "*="
				} else if data[i+1] == '*' {
					i += 2
					token.tokenType = TkSqrInc
					token.sdata = "**"
				} else {
					i++
					token.tokenType = TkMul
					token.sdata = "*"
				}

			case "%":
				if data[i+1] == '=' {
					i += 2
					token.tokenType = TkRemAssign
					token.sdata = "%="
				} else {
					i++
					token.tokenType = TkRem
					token.sdata = "%"
				}

			case "&":
				if data[i+1] == '^' && data[i+2] == '=' {
					i += 3
					token.tokenType = TkAndNotAssign
					token.sdata = "&^="
				} else if data[i+1] == '=' {
					i += 2
					token.tokenType = TkAndAssign
					token.sdata = "&="
				} else if data[i+1] == '^' {
					i += 2
					token.tokenType = TkAndNot
					token.sdata = "&^"
				} else if data[i+1] == '&' {
					i += 2
					token.tokenType = TkLand
					token.sdata = "&&"
				} else {
					i++
					token.tokenType = TkAnd
					token.sdata = "&"
				}

			case "|":
				if data[i+1] == '=' {
					i += 2
					token.tokenType = TkOrAssign
					token.sdata = "|="
				} else if data[i+1] == '|' {
					i += 2
					token.tokenType = TkLor
					token.sdata = "||"
				} else {
					i++
					token.tokenType = TkOr
					token.sdata = "|"
				}

			case "^":
				if data[i+1] == '=' {
					i += 2
					token.tokenType = TkXorAssign
					token.sdata = "^="
				} else {
					i++
					token.tokenType = TkXor
					token.sdata = "^"
				}

			case "+":
				if data[i+1] == '=' {
					i += 2
					token.tokenType = TkAddAssign
					token.sdata = "+="
				} else if data[i+1] == '+' {
					i += 2
					token.tokenType = TkInc
					token.sdata = "++"
				} else {
					i++
					token.tokenType = TkAdd
					token.sdata = "+"
				}

			case "-":
				if data[i+1] == '=' {
					i += 2
					token.tokenType = TkSubAssign
					token.sdata = "-="
				} else if data[i+1] == '-' {
					i += 2
					token.tokenType = TkDec
					token.sdata = "--"
				} else {
					i++
					token.tokenType = TkSub
					token.sdata = "-"
				}

			case "'":
				i++
				c = data[i]

				token.tokenType = TkChar

				if c == '\\' {
					i++
					c = data[i]
					token.sdata = "\\" + string(c)
				} else {
					token.sdata += string(c)
				}

				i++
				c = data[i]

				if c != '\'' {
					errs = append(errs, errors.New("Character literal has multiple characters: "+token.toString()))
				}

				i++

			case "<":
				if data[i+1] == '<' && data[i+2] == '=' {
					i += 3
					token.tokenType = TkShlAssign
					token.sdata = "<<="
				} else if data[i+1] == '<' {
					i += 2
					token.tokenType = TkShl
					token.sdata = "<<"
				} else if data[i+1] == '-' {
					i += 2
					token.tokenType = TkArrow
					token.sdata = "<-"
				} else if data[i+1] == '=' {
					i += 2
					token.tokenType = TkLeq
					token.sdata = "<="
				} else {
					i++
					token.tokenType = TkLss
					token.sdata = "<"
				}

			case ">":
				if data[i+1] == '>' && data[i+2] == '=' {
					i += 3
					token.tokenType = TkShrAssign
					token.sdata = ">>="
				} else if data[i+1] == '>' {
					i += 2
					token.tokenType = TkShr
					token.sdata = ">>"
				} else if data[i+1] == '=' {
					i += 2
					token.tokenType = TkGeq
					token.sdata = ">="
				} else {
					i++
					token.tokenType = TkGtr
					token.sdata = ">"
				}

			case "=":
				if data[i+1] == '=' {
					i += 2
					token.tokenType = TkEql
					token.sdata = "=="
				} else {
					i++
					token.tokenType = TkAssign
					token.sdata = "="
				}

			case "!":
				if data[i+1] == '=' {
					i += 2
					token.tokenType = TkNeq
					token.sdata = "!="
				} else {
					i++
					token.tokenType = TkNot
					token.sdata = "!"
				}

			case "(":
				i++
				token.tokenType = TkLParen
				token.sdata = "("

			case ")":
				i++
				token.tokenType = TkRParen
				token.sdata = ")"

			case "[":
				i++
				token.tokenType = TkLBrack
				token.sdata = "["

			case "]":
				i++
				token.tokenType = TkRBrack
				token.sdata = "]"

			case "{":
				i++
				token.tokenType = TkLBrace
				token.sdata = "{"

			case "}":
				i++
				token.tokenType = TkRBrace
				token.sdata = "}"

			case ",":
				i++
				token.tokenType = TkComma
				token.sdata = ","

			default:
				errs = append(errs, errors.New("Unrecognized character: "+string(c)))
			}
		}

		if token.tokenType != TkIllegal && token.tokenType != TkComment {
			tokens = append(tokens, token)
		}
	}

	tokens = append(tokens, Token{1, ln + 1, "__eof__", 0, 0.0})

	return tokens, errs
}
