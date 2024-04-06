package token

import (
	"fmt"
)

// Pos 类型定义，用于存储位置信息
type Pos struct {
	Offset int
	Line   int
	Col    int
}

func (pos Pos) String() string {
	return fmt.Sprintf("%d:%d:%d", pos.Offset, pos.Line, pos.Col)
}
