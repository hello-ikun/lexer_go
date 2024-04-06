// 代码参考go语言官方代码
package token

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

// TokenType是Go编程语言的词法标记集合。
type TokenType int

// 标记列表。
const (
	// 特殊标记
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	// 标识符和基本类型字面量（这些标记代表字面量的类别）
	literal_beg
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"
	BLEAK

	literal_end

	// 操作符和分隔符
	operator_beg
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :
	operator_end

	keyword_beg
	// 关键字
	BREAK
	CASE
	CHAN
	CONST
	CONTINUE

	DEFAULT
	DEFER
	ELSE
	FALLTHROUGH
	FOR

	FUNC
	GO
	GOTO
	IF
	IMPORT

	INTERFACE
	MAP
	PACKAGE
	RANGE
	RETURN

	SELECT
	STRUCT
	SWITCH
	TYPE
	VAR
	keyword_end

	additional_beg
	// 额外的标记，以临时方式处理
	TILDE
	additional_end
)

var TokenTypes = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	IMAG:   "IMAG",
	CHAR:   "CHAR",
	STRING: "STRING",
	BLEAK:  "_",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND:     "&",
	OR:      "|",
	XOR:     "^",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",

	LAND:  "&&",
	LOR:   "||",
	ARROW: "<-",
	INC:   "++",
	DEC:   "--",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ:      "!=",
	LEQ:      "<=",
	GEQ:      ">=",
	DEFINE:   ":=",
	ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	BREAK:    "break",
	CASE:     "case",
	CHAN:     "chan",
	CONST:    "const",
	CONTINUE: "continue",

	DEFAULT:     "default",
	DEFER:       "defer",
	ELSE:        "else",
	FALLTHROUGH: "fallthrough",
	FOR:         "for",

	FUNC:   "func",
	GO:     "go",
	GOTO:   "goto",
	IF:     "if",
	IMPORT: "import",

	INTERFACE: "interface",
	MAP:       "map",
	PACKAGE:   "package",
	RANGE:     "range",
	RETURN:    "return",

	SELECT: "select",
	STRUCT: "struct",
	SWITCH: "switch",
	TYPE:   "type",
	VAR:    "var",

	TILDE: "~",
}

// String返回与标记tok对应的字符串。
func (tok TokenType) String() string {
	s := ""
	if 0 <= tok && tok < TokenType(len(TokenTypes)) {
		s = TokenTypes[tok]
	}
	if s == "" {
		s = "TokenType(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]TokenType
var opMap map[string]TokenType

// 建立keywords映射关系
func init() {
	keywords = make(map[string]TokenType, keyword_end-(keyword_beg+1))
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[TokenTypes[i]] = i
	}
	opMap = make(map[string]TokenType, operator_end-(operator_beg+1))
	for i := operator_beg + 1; i < operator_end; i++ {
		keywords[TokenTypes[i]] = i
	}
}

// Lookup将标识符映射到其关键字标记或[IDENT]（如果不是关键字）。
func Lookup(ident string) TokenType {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	if tok, is_op := opMap[ident]; is_op {
		return tok
	}
	return IDENT
}

// IsLiteral对应于标识符和基本类型字面量的标记返回true；否则返回false。
func (tok TokenType) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator对应于操作符和分隔符的标记返回true 否则返回false。
func (tok TokenType) IsOperator() bool {
	return (operator_beg < tok && tok < operator_end) || tok == TILDE
}

// IsKeyword对应于关键字的标记返回true； 否则返回false
func (tok TokenType) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

// IsExported报告名称是否以大写字母开头。
func IsExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

// IsKeyword报告名称是否是Go关键字，如“func”或“return”。
func IsKeyword(name string) bool {
	_, ok := keywords[name]
	return ok
}

// IsIdentifier报告名称是否是Go标识符，即由字母、数字和下划线组成的非空字符串，
// 其中第一个字符不是数字。关键字不是标识符。
func IsIdentifier(name string) bool {
	if name == "" || IsKeyword(name) {
		return false
	}
	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return true
}
