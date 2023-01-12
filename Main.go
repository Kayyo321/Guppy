package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Println("gup <flags> <files>")
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	filename := os.Args[1]
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading: \""+filename+"\":", err)
		os.Exit(1)
	}

	data := string(bytes)

	tokens, errs := lex(data)
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println("<!>", err)
		}

		fmt.Println("Compiler exited from the lexer.")
		os.Exit(1)
	}

	var p Parser
	p.set(tokens)

	ast, errs := p.parse()
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println("<!>", err)
		}

		fmt.Println("Compiler exited from the parser.")
		os.Exit(1)
	}

	fmt.Println("AST {\n" + ast.toString() + "\n}")
}
