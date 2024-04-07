package token

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestToken(t *testing.T) {
	fp, err := os.Open("2.txt")
	if err != nil {
		panic(err)
	}
	buf, err := io.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	sourceCode := string(buf)
	// 初始化一个缓冲区用于存储多行注释
	var commentBuffer bytes.Buffer

	// 遍历每个字节
	for i := 0; i < len(sourceCode); i++ {
		// 检查是否遇到了多行注释的起始
		if i+1 < len(sourceCode) && sourceCode[i] == '/' && sourceCode[i+1] == '*' {
			// 将多行注释的起始添加到缓冲区

			commentBuffer.WriteByte(sourceCode[i])
			commentBuffer.WriteByte(sourceCode[i+1])
			// 继续读取下一个字节
			i += 2

			// 继续读取直到找到多行注释的结束
			for ; i < len(sourceCode); i++ {
				commentBuffer.WriteByte(sourceCode[i])
				if i+1 < len(sourceCode) && sourceCode[i] == '*' && sourceCode[i+1] == '/' {
					// 多行注释的结束，跳出循环
					commentBuffer.WriteByte(sourceCode[i+1])
					i++
					break
				}
			}

			// 打印拼接后的多行注释
			ans := commentBuffer.String()
			fmt.Println(ans)
			// 清空缓冲区
			commentBuffer.Reset()
		}
	}
}
