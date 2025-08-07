package test

import (
	"fmt"
	"testing"

	"github.com/szks-repo/gosmarty/parser"
	"github.com/szks-repo/gosmarty/lexer"
)

func Test_Day1(t *testing.T) {
	input := `Hello {$name}!`

	l := lexer.New(input)

	p := parser.New(l)

	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) != 0 {
		fmt.Println("Woops! We ran into some monkey business here:")
		for _, msg := range errors {
			fmt.Printf("\t- %s\n", msg)
		}
		return
	}

	// 5. 結果（AST）を表示
	// ast.ProgramのString()メソッドが呼ばれる
	fmt.Println("--- AST ---")
	fmt.Println(program.String())
}
