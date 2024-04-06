package token

import (
	"fmt"
	"strings"
)

// Token 类型定义
type Token struct {
	TokenType TokenType
	Value     interface{}
	Pos       Pos
}

// Scanner 类定义
type Scanner struct {
	InputString string
	Position    int
	Line        int
	Col         int
}

// 构造函数
func NewScanner(inputString string) *Scanner {
	return &Scanner{
		InputString: inputString,
		Position:    0,
		Line:        1,
		Col:         1,
	}
}

// 词法分析函数
func (s *Scanner) Scan() []Token {
	var tokens []Token

	for s.Position < len(s.InputString) {
		char := s.InputString[s.Position]

		if char == ' ' || char == '\t' || char == '\r' {
			s.updatePos()
		} else if char == '\n' {
			s.updatePos()
		} else if char == '"' || char == '`' || char == '\'' {
			tokens = append(tokens, s.extractStringAndByte())
		} else if s.isAlpha(char) || char == '_' {
			token := s.extractIdentifier()
			tokens = append(tokens, token)
		} else if s.isDigit(char) || (char == '-' && s.isDigit(s.InputString[s.Position+1])) {
			tokens = append(tokens, s.extractNumber())
		} else if strings.Contains("+-*%=()&|^<>!.:;/{}|[]^<>\\,", string(char)) { // 包含字符
			token := s.extractOperator()
			tokens = append(tokens, token)
		} else {
			panic(fmt.Sprintf("错误/未知字符 '%c' 出现在 %d", char, s.Position))
		}
	}

	return tokens
}

// 词法分析函数
func (s *Scanner) FScan() {
	// 输出标题
	fmt.Printf("%-18s %-15s %s\n", "Position", "Type", "Value")
	var tok Token
	for s.Position < len(s.InputString) {
		char := s.InputString[s.Position]

		if char == ' ' || char == '\t' || char == '\r' {
			s.updatePos()
			continue
		} else if char == '\n' {
			s.updatePos()
			continue
		} else if char == '"' || char == '`' || char == '\'' {
			tok = s.extractStringAndByte()
		} else if s.isAlpha(char) || char == '_' {
			tok = s.extractIdentifier()

		} else if s.isDigit(char) {
			tok = s.extractNumber()
		} else if strings.Contains("+-*%=()&|^<>!.:;/{}|[]^<>\\,", string(char)) { // 包含字符
			tok = s.extractOperator()
		} else {
			panic(fmt.Sprintf("错误/未知字符 '%c' 出现在 %d", char, s.Position))
		}
		fmt.Printf("%-6d:%-6d:%-6d %-15s %v\n", s.Position, tok.Pos.Line, tok.Pos.Col, tok.TokenType, tok.Value)
	}

}

// 提取分割符号符
func (s *Scanner) extractOperator() Token {
	startPosition := s.Position
	startLine := s.Line
	startCol := s.Col
	operator := string(s.InputString[s.Position])
	s.updatePos()

	switch operator {
	case "+":
		// ++
		if s.InputString[s.Position] == '+' {
			operator += "+"
			s.updatePos()
		}
		// +=
		if s.InputString[s.Position] == '=' {
			operator += "="
			s.updatePos()
		}
	case "-":
		// --
		if s.InputString[s.Position] == '-' {
			operator += "-"
			s.updatePos()
		}
		// -=
		if s.InputString[s.Position] == '=' {
			operator += "="
			s.updatePos()
		}
		// 负数的情况 在这里我们是不需要进行考虑的
	case "*":
		// *=
		if s.InputString[s.Position] == '=' {
			operator += "="
			s.updatePos()
		}
	case "/":
		// 单行注释
		if s.InputString[s.Position] == '/' {
			return s.extractSingleLineComment()
		}
		// 多行注释
		if s.InputString[s.Position] == '*' {
			return s.extractMultiLineComment()
		}
		// /=
		if s.InputString[s.Position] == '=' {
			operator += "="
			s.updatePos()
		}
	case "%":
		// %=
		if s.InputString[s.Position] == '=' {
			operator += "="
			s.updatePos()
		}
	case "&":
		// &^
		if s.InputString[s.Position] == '^' {
			operator += "^"
			s.updatePos()
		}
		// &=
		if s.InputString[s.Position] == '=' {
			operator += "="
			s.updatePos()
		}
	case "|":
		// |=
		if s.InputString[s.Position] == '=' {
			operator += "="
			s.updatePos()
		}
	case "^":
		// ^=
		if s.InputString[s.Position] == '=' {
			operator += "="
			s.updatePos()
		}
	case "<":
		// <<
		if s.InputString[s.Position] == '<' {
			operator += "<"
			s.updatePos()
			// <<=
			if s.InputString[s.Position] == '=' {
				operator += "="
				s.updatePos()
			}
		}
		// <=
		if s.InputString[s.Position] == '=' {
			operator += "="
			s.updatePos()
		}
	case ">":
		// >>
		if s.InputString[s.Position] == '>' {
			operator += ">"
			s.updatePos()
			// >>=
			if s.InputString[s.Position] == '=' {
				operator += "="
				s.updatePos()
			}
		}
		// >=
		if s.InputString[s.Position] == '=' {
			operator += "="
			s.updatePos()
		}
	}

	pos := Pos{Offset: startPosition, Line: startLine, Col: startCol}
	value := s.InputString[startPosition:s.Position]

	tok := Lookup(value) // 检查是不是关键字
	return Token{TokenType: tok, Value: operator, Pos: pos}
}

// 更新位置信息函数
func (s *Scanner) updatePos() {
	if s.InputString[s.Position] == '\n' {
		s.Line++
		s.Col = 1
	} else {
		s.Col++
	}
	s.Position++
}

// 判断字符是否为字母
func (s *Scanner) isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

// 判断字符是否为数字
func (s *Scanner) isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

// 提取标识符
func (s *Scanner) extractIdentifier() Token {
	startPosition := s.Position
	startLine := s.Line
	startCol := s.Col

	for s.Position < len(s.InputString) && (s.isAlpha(s.InputString[s.Position]) || s.isDigit(s.InputString[s.Position]) || s.InputString[s.Position] == '_') {
		s.updatePos()
	}

	value := s.InputString[startPosition:s.Position]

	tok := Lookup(value) // 检查是不是关键字
	token := Token{TokenType: tok, Value: value, Pos: Pos{Offset: startPosition, Line: startLine, Col: startCol}}

	return token
}

// 提取数字
func (s *Scanner) extractNumber() Token {
	startPosition := s.Position
	startLine := s.Line
	startCol := s.Col

	var tokenType TokenType
	// 提取整数部分
	for s.Position < len(s.InputString) && s.isDigit(s.InputString[s.Position]) {
		s.updatePos()
		tokenType = INT
	}

	// 提取小数部分
	if s.Position < len(s.InputString) && s.InputString[s.Position] == '.' {
		s.updatePos()
		for s.Position < len(s.InputString) && s.isDigit(s.InputString[s.Position]) {
			s.updatePos()
		}
		tokenType = FLOAT
	}

	// 提取科学计数法部分
	if s.Position < len(s.InputString) && (s.InputString[s.Position] == 'e' || s.InputString[s.Position] == 'E') {
		s.updatePos()
		if s.Position < len(s.InputString) && (s.InputString[s.Position] == '+' || s.InputString[s.Position] == '-') {
			s.updatePos()
		}
		for s.Position < len(s.InputString) && s.isDigit(s.InputString[s.Position]) {
			s.updatePos()
		}
		tokenType = FLOAT
	}
	if s.Position+1 < len(s.InputString) && (s.InputString[s.Position] == 'i') {
		s.updatePos()
		tokenType = IMAG
	}
	value := s.InputString[startPosition:s.Position]
	token := Token{TokenType: tokenType, Value: value, Pos: Pos{Offset: startPosition, Line: startLine, Col: startCol}}

	return token
}
func (s *Scanner) extractStringAndByte() Token {
	Pre := s.InputString[s.Position]
	StartLine := s.Line
	StartCol := s.Col
	StartPosition := s.Position
	s.Position++
	escape := false
	value := string(Pre)

	for s.Position < len(s.InputString) && (s.InputString[s.Position] != Pre || escape) {
		if escape {
			value += "\\" + string(s.InputString[s.Position])
			escape = false
		} else {
			if s.InputString[s.Position] == '\\' {
				escape = true
			} else {
				value += string(s.InputString[s.Position])
			}
		}
		s.updatePos()
	}

	if s.Position == len(s.InputString) {
		panic("字符串解析存在错误")
	}

	s.updatePos()

	// 你需要根据你的实际代码来定义 Pos 类型的结构体和 UpdatePos() 方法

	Pos := Pos{StartPosition, StartLine, StartCol}
	tokenType := STRING
	if Pre == '\'' {
		tokenType = CHAR
	}

	return Token{tokenType, value + string(Pre), Pos}
}

// 提取单行注释
func (s *Scanner) extractSingleLineComment() Token {
	startPosition := s.Position - 1
	startLine := s.Line
	startCol := s.Col - 1

	for s.Position < len(s.InputString) && s.InputString[s.Position] != '\n' {
		s.updatePos()
	}
	value := strings.TrimSpace(s.InputString[startPosition:s.Position])
	pos := Pos{Offset: startPosition, Line: startLine, Col: startCol}
	return Token{TokenType: COMMENT, Value: value, Pos: pos}
}

// 提取多行注释
func (s *Scanner) extractMultiLineComment() Token {
	startPosition := s.Position - 1
	startLine := s.Line
	startCol := s.Col - 1

	for s.Position < len(s.InputString)-1 && s.InputString[s.Position:s.Position+2] != "*/" {
		s.updatePos()
	}
	if s.Position == len(s.InputString)-1 {
		panic("多行注释出现错误")
	}
	// 进行转义处理
	s.updatePos()
	s.updatePos()
	value := strings.ReplaceAll(s.InputString[startPosition:s.Position], "\n", "\\n")
	pos := Pos{Offset: startPosition, Line: startLine, Col: startCol}
	return Token{TokenType: COMMENT, Value: value, Pos: pos}
}

// 测试函数
func TestLexer(inputString string) {
	scanner := NewScanner(inputString)
	tokens := scanner.Scan()

	// 输出标题
	fmt.Printf("%-15s %-15s %s\n", "Position", "Type", "Value")

	// 输出词法分析结果
	for _, token := range tokens {
		position := fmt.Sprintf("%d:%d", token.Pos.Line, token.Pos.Col)
		fmt.Printf("%-15s %-15s %v\n", position, token.TokenType, token.Value)
	}
}

// 测试函数
func TestLexerFormat(inputString string) {
	scanner := NewScanner(inputString)
	scanner.FScan()

}
