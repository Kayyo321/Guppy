package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scnr := bufio.NewScanner(os.Stdin)

	var data string

	for {
		fmt.Print("comp (\"~run\" to run file) $ ")
		scnr.Scan()

		txt := scnr.Text()
		if txt == "~run" {
			break
		} else {
			data += txt
		}
	}

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
