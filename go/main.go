package main

import (
	"fmt"
	"io"
	"lexer/token"
	"os"
)

/*
这个是多行注释
*/
func main() {
	x := 12 + 2.3131i

	name := `
	eksfnskf
	`
	fmt.Println(x, name)
	fp, err := os.Open("main.go")
	if err != nil {
		panic(err)
	}
	buf, err := io.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	inputString := string(buf)
	token.TestLexerFormat(inputString)
}
