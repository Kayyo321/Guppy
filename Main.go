package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Println("Usage: gup <interaction> <flags> <files>")
	fmt.Println()
	fmt.Println("<interaction>'s")
	fmt.Println("    ├──> build : compile the following files into an executable.")
	fmt.Println("    ├──> run   : compile the files, and run it after built.")
	fmt.Println("    └──> edit  : edit a virtual file line by line, then save and or run it.")
	fmt.Println()
	fmt.Println("<flags>")
	fmt.Println("    ├──> '-o'   : build the files, and rename the output file.")
	fmt.Println("    ├──> '-asm' : build the files to x86 assembly; instead of '.exe' || '.o'.")
	fmt.Println("    └──> <!> TODO: think of more flags <!>")
	fmt.Println()
	fmt.Println("<files>")
	fmt.Println("    ├──> Anything that ends with '.gpy' || '.guppy' .")
	fmt.Println("    ├──> Support absolute && relative file path.")
	fmt.Println("    └──> (by default, output name is \"a.exe\", not the first file's name.")
	os.Exit(0)
}

func build(files []string, flags map[string]string) string {
	// var obj []string
	for _, filename := range files {
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

		// obj = append(obj, outputPath)
	}

	// After compilation, link all objects.
	// Return filepath to exe.

	return ""
}

func run(files []string, flags map[string]string) {
	o := build(files, flags)
	fmt.Println(o)
	// Run the exe.
}

func parseArgs() (map[string]string, []string) {
	flags := make(map[string]string)
	var files []string

	for i, arg := range os.Args {
		if arg[0] == '-' {
			flags[arg] = os.Args[i+1]
		} else {
			fileExt := filepath.Ext(arg)
			if fileExt == ".gpy" || fileExt == ".guppy" {
				files = append(files, arg)
			}
		}
	}

	return flags, files
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	// TODO: impliment flags.
	flags, files := parseArgs()
	fmt.Println(flags)

	switch os.Args[1] {
	case "build":
		build(files, flags)

	case "run":
		run(files, flags)

	case "edit":
		break

	default:
		fmt.Println("Unknown interaction:", os.Args[1])
		usage()
		os.Exit(1)
	}
}
