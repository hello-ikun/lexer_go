# GO语言词法分析器
>使用go和python分别实现

## Go部分
>已经写完了,可能存在部分未知问题 | 简单做一个描述
### Token和Pos结构
```go
# 定义 token 类型
// TokenType是Go编程语言的词法标记集合。
type TokenType int

// 标记列表。
const (
	// 特殊标记
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	// 标识符和基本类型字面量（这些标记代表字面量的类别）
	//此处省略
)
```
```go
// Pos 类型定义，用于存储位置信息
type Pos struct {
	Offset int
	Line   int
	Col    int
}

func (pos Pos) String() string {
	return fmt.Sprintf("%d:%d:%d", pos.Offset, pos.Line, pos.Col)
}
```
### Scanner结构|主要部分
```go
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
```
>第一步:提取标识符
```go
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
```
>第二步:对于操作符/分割符号进行解析|注意这个符号有时候有歧义
```go
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
```
>第三步:解析数字|可以有一个更简单的做法(但我们不采用)
```go
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
```
>第四步:解析一下字符类型|`"'这三种 |注意转义处理
```go
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
	Pos := Pos{StartPosition, StartLine, StartCol}
	tokenType := STRING
	if Pre == '\'' {
		tokenType = CHAR
	}
	value = strings.ReplaceAll(strings.TrimSpace(value), "\r\n", "\\n")
	return Token{tokenType, value + string(Pre), Pos}
}
```
> 第五步:解析注释信息 | 单行和多行注释
```go
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

// extractMultiLineComment 提取多行注释
func (s *Scanner) extractMultiLineComment() Token {
	startPosition := max(s.Position-1, 0)
	startLine := s.Line
	startCol := s.Col - 1
	for s.Position+2 < len(s.InputString)-1 && s.InputString[s.Position:s.Position+2] != "*/" {
		s.updatePos()
	}

	if s.Position == len(s.InputString)-1 {
		panic("多行注释出现错误")
	}

	// 进行转义处理
	s.updatePos()
	s.updatePos()
	value := strings.ReplaceAll(s.InputString[startPosition:s.Position], "\r\n", "\\n")
	pos := Pos{Offset: startPosition, Line: startLine, Col: startCol}
	return Token{TokenType: COMMENT, Value: value, Pos: pos}
}
```
>补充一下 我们需要更新位置信息以及辅助函数
```go
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
```
>最后来写调用scan函数
```go
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
```